package web

import (
	"net/http"
	"os"

	"github.com/chuckha/kubeyaml.com/backend/internal/shared/logging"
)

const (
	defaultAddr    = ":9000"
	defaultDevMode = false
)

type service interface {
	Validate([]byte) error
}

type logger interface {
	Infof(a string, args ...interface{})
}

type Server struct {
	svr *http.Server
	dev bool
	svc service
	log logger
}

type ServerOption func(s *Server)

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.svr.Addr = addr
	}
}

func WithDevMode(dev bool) ServerOption {
	return func(s *Server) {
		s.dev = dev
	}
}

func NewServer(svc service, opts ...ServerOption) *Server {
	s := &Server{
		dev: defaultDevMode,
		log: &logging.Log{Writer: os.Stdout},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", s.corsForDev(s.validate))
	mux.HandleFunc("/favicon.ico", s.corsForDev(s.favicon))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", s.main)
	svr := &http.Server{
		Addr:    defaultAddr,
		Handler: mux,
	}
	s.svr = svr
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Server) Run() {
	s.log.Infof("Serving web traffic on %s\n", s.svr.Addr)
	s.log.Infof("Development mode is %v\n", devMode(s.dev))
	s.svr.ListenAndServe()
}

func devMode(dev bool) string {
	if dev {
		return "on"
	}
	return "off"
}
