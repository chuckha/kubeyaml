package web

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/chuckha/kubeyaml.com/backend/internal/kubernetes"
)

func (s *Server) main(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("internal/adapters/web/templates/main.go.tmpl"))
	tmpl.Execute(w, nil)
}

func (s *Server) favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.png")
}

func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	s.logRequest("validate", r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Infof("error reading body: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// data is posted with plain HTML so we get `data=url+encoded+yaml&key=value...`
	v, err := url.ParseQuery(string(b))
	if err != nil {
		s.log.Infof("error parsing value string: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Extract the yaml and load it
	data := v.Get("data")
	if len(data) == 0 {
		// Ignore empty requests
		return
	}
	if err := s.svc.Validate([]byte(data)); err != nil {
		switch err.(type) {
		case *kubernetes.RequiredKeyNotFoundError:
		case *kubernetes.YamlPathError:
		default:
			s.log.Infof("error loading body with non user error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		errs := make(map[string][]error)
		// TODO: use the incoming versions instead of the hardcoded ons
		// for _, v := range s.validators {
		// 	errs[v.Version()] = []error{err}
		// }

		out, err := json.Marshal(errs)
		if err != nil {
			s.log.Infof("error marshalling errors: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if _, err := w.Write(out); err != nil {
			s.log.Infof("error writing response body: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
}
