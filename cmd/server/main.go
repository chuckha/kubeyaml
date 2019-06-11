package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chuckha/kubeyaml/internal/kubernetes"
)

const (
	templateSuffix = ".template.html"
)

type ServerArgs struct {
	Port        string
	Development bool
}

func main() {
	// Set up the server args
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	sa := &ServerArgs{}
	fs.StringVar(&sa.Port, "port", "9000", "the port for the server to listen on")
	fs.BoolVar(&sa.Development, "dev", false, "enable certain features when developing locally")

	// Parse flags
	fs.Parse(os.Args[1:])

	versions := []string{"1.11", "1.12", "1.13", "1.14"}
	validators := make([]validator, len(versions))
	for i, version := range versions {
		resolver, err := kubernetes.NewResolver(version)
		if err != nil {
			fmt.Printf("failed to get a resolver for version %q: %v", version, err)
			os.Exit(1)
		}
		validators[i] = kubernetes.NewValidator(resolver)
	}

	loader := kubernetes.NewLoader()

	// this is a bad optimization. This is essentially sharing the group finder across
	// all versions of kubernetes apis. It's entirely possible api versions have
	// different namespaces.
	// TODO associate this with the resolver and expose through the validator.
	gf := kubernetes.NewAPIKeyer("io.k8s.api", ".k8s.io")

	s := &server{
		logger:     &log{os.Stdout},
		validators: validators,
		loader:     loader,
		finder:     gf,
		dev:        sa.Development,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", s.validate)
	mux.HandleFunc("/favicon.ico", s.favicon)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	fmt.Printf("listening on port :%s\n", sa.Port)
	if sa.Development {
		fmt.Println("development mode enabled")
	}
	http.ListenAndServe(":"+sa.Port, mux)
}

type logger interface {
	Debugf(string, ...interface{}) error
	Infof(string, ...interface{}) error
}

type log struct {
	writer io.Writer
}

func (l *log) Debugf(format string, args ...interface{}) error {
	wrappedFormat := fmt.Sprintf("[DEBUG]: %s", format)
	_, err := fmt.Fprintf(l.writer, wrappedFormat, args...)
	return err
}
func (l *log) Infof(format string, args ...interface{}) error {
	wrappedFormat := fmt.Sprintf("[INFO]: %s", format)
	_, err := fmt.Fprintf(l.writer, wrappedFormat, args...)
	return err
}

type validator interface {
	Validate(map[interface{}]interface{}, *kubernetes.Schema) []error
	Resolve(string) (*kubernetes.Schema, error)
	Version() string
}

type loader interface {
	Load(io.Reader) (*kubernetes.Input, error)
}
type groupFinder interface {
	APIKey(string, string) string
}

type server struct {
	logger
	validators []validator
	loader     loader
	finder     groupFinder
	dev        bool
}

func (s *server) favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func (s *server) validate(w http.ResponseWriter, r *http.Request) {
	// Enable CORS in dev mode
	if s.dev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	s.logRequest("validate", r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Infof("error reading body: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// data is posted with plain HTML so we get `data=url+encoded+yaml&key=value...`
	v, err := url.ParseQuery(string(b))
	if err != nil {
		s.logger.Infof("error parsing value string: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Extract the yaml and load it
	data := strings.NewReader(v.Get("data"))
	i, err := s.loader.Load(data)
	if err != nil {
		switch err.(type) {
		case *kubernetes.RequiredKeyNotFoundError:
		case *kubernetes.YamlPathError:
		default:
			s.logger.Infof("error loading body with non user error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		errs := make(map[string][]error)
		for _, v := range s.validators {
			errs[v.Version()] = []error{err}
		}

		out, err := json.Marshal(errs)
		if err != nil {
			s.logger.Infof("error marshalling errors: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if _, err := w.Write(out); err != nil {
			s.logger.Infof("error writing response body: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	defer r.Body.Close()

	// Validate each version
	errs := make(map[string][]error)
	for _, v := range s.validators {

		// Lookup the api group version to get started and ensure the kind is valid
		// TODO: Consider making Resolve take a version and a kind so it can produce a better error message
		schema, err := v.Resolve(s.finder.APIKey(i.APIVersion, i.Kind))
		if err != nil {
			errs[v.Version()] = []error{err}
			continue
		}

		errs[v.Version()] = v.Validate(i.Data, schema)
	}

	out, err := json.Marshal(errs)
	if err != nil {
		s.logger.Infof("error marshalling errors: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(out); err != nil {
		s.logger.Infof("error writing response body: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) logRequest(method string, r *http.Request) {
	s.logger.Infof("[%s]: %s %s\n", method, r.Method, r.URL.Path)
}
