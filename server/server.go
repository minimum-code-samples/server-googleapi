package server

import (
	"github.com/gorilla/mux"
)

// Server holds the shared dependencies for the Web server.
type Server struct {
	// The interface the server listens on.
	Interface string
	// The port that the server listens on.
	Port string
	// The port that the server listens on for TLS connections.
	PortTLS string
	// Whether the server has been bootstrapped successfully.
	Ready bool
	// The router.
	Router *mux.Router
	// The Config instance.
	config Config
}

// NewServer creates a new instance of the server.
func NewServer(cfg Config) *Server {
	return &Server{
		Interface: cfg.Interface,
		Port:      cfg.Port,
		Ready:     false,
		Router:    nil,
		config:    cfg,
	}
}

// MakeRouter creates the router ready with the handlers mapped to all the paths.
func (s *Server) MakeRouter(useMiddleware bool) {
	s.Router = mux.NewRouter()
	// if (useMiddleware) {
	// 	s.Router.Use(middleware.HSTS)
	// }
	// Create subrouter for API endpoints.
	// api := s.Router.PathPrefix("/v1").Subrouter()()
	// routeAPI(api, s, useMiddleware)
	routePages(s, useMiddleware)
	routeStatic(s)
}
