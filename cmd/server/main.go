package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/chuckha/kube-validate/internal/kubernetes"
)

const (
	templateSuffix = ".template.html"
)

func main() {
	versions := []string{"1.8", "1.9", "1.10", "1.11", "1.12"}
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

	t, err := loadTemplates()
	if err != nil {
		fmt.Printf("failed to load templates: %v\n", err)
		os.Exit(1)
	}
	s := &server{
		templates:  t,
		logger:     &log{os.Stdout},
		validators: validators,
		loader:     loader,
		finder:     gf,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", s.validate)
	mux.HandleFunc("/favicon.ico", s.favicon)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", s.index)
	fmt.Println("listening on :9000")
	http.ListenAndServe(":9000", mux)
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
	Validate(map[string]interface{}, *kubernetes.Schema, []string) []error
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
	templates map[string]*template.Template
	logger
	validators []validator
	loader     loader
	finder     groupFinder
}

func (s *server) favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	// TODO generate html based on versions maybe
	s.logRequest("index", r)
	if err := s.templates["index"].Execute(w, nil); err != nil {
		http.Error(w, "failed to execute index template", http.StatusInternalServerError)
		return
	}
}

func (s *server) validate(w http.ResponseWriter, r *http.Request) {
	s.logRequest("validate", r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Infof("error loading body: %v\n", err)
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
		s.logger.Infof("error loading body: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate each version
	errs := make(map[string][]error)
	for _, v := range s.validators {

		// Lookup the api group version to get started and ensure the kind is valid
		schema, err := v.Resolve(s.finder.APIKey(i.APIVersion, i.Kind))
		if err != nil {
			errs[v.Version()] = []error{err}
			continue
		}

		errs[v.Version()] = v.Validate(i.Data, schema, []string{})
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

func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	filepath.Walk("templates", func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walkfn received an error: %v", err)
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(p, templateSuffix) {
			tmpl, err := template.ParseFiles(p)
			if err != nil {
				return fmt.Errorf("failed to parse template %q: %v", p, err)
			}
			templates[path.Base(strings.TrimSuffix(p, templateSuffix))] = tmpl
		}

		return nil
	})
	return templates, nil
}
