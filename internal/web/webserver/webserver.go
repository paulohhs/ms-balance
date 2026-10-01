package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort string
}

func NewWebServer(webServerPort string) *WebServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &WebServer{
		Router:        router,
		WebServerPort: webServerPort,
	}
}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	s.Router.Method(method, path, handler)
}

func (s *WebServer) Start() error {
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
