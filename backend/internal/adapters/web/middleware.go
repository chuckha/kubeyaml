package web

import "net/http"

func (s *Server) logRequest(method string, r *http.Request) {
	s.log.Infof("[%s]: %s %s\n", method, r.Method, r.URL.Path)
}

// corsForDev enables CORS when the server is running in development mode
func (s *Server) corsForDev(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.dev {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		f(w, r)
	}
}
