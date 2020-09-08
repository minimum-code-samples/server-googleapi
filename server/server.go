package server

import (
	"net/http"
	"server-googleapi/google"
	"server-googleapi/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

// Server holds the shared dependencies for the Web server.
type Server struct {
	// The interface the server listens on.
	Interface string
	// The port that the server listens on.
	Port string
	// The port that the server listens on for TLS connections.
	PortTLS string
	// The router.
	Router *mux.Router
	// The session store.
	SessStore sessions.Store
	// The token for the admin.
	TokenAdmin *oauth2.Token
	// The Config instance.
	config Config
}

// NewServer creates a new instance of the server.
func NewServer(cfg Config, store sessions.Store) *Server {
	return &Server{
		Interface: cfg.Interface,
		Port:      cfg.Port,
		Router:    nil,
		SessStore: store,
		config:    cfg,
	}
}

// Auth stores user information into session cookies.
func (s *Server) Auth(details *google.UserDetail, w http.ResponseWriter, r *http.Request) bool {
	sess := s.RedirectIfUnauth(w, r)
	if sess == nil {
		return false
	}
	sess.Values[model.SessName] = details.Name
	if e := sess.Save(r, w); e != nil {
		s.RedirectError(w, r, model.ErrorSessionSave)
		return false
	}
	return true
}

// RedirectError is a helper method to reduce boilerplate code.
func (s *Server) RedirectError(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, model.PathError+"?msg="+msg, http.StatusTemporaryRedirect)
	// If StatusInternalServerError, shows an ugly "Internal Server Error" page.
}

// RedirectIfUnauth redirects to the error page if session is not valid.
//
// Returns nil if redirection occurs, the session object otherwise.
func (s *Server) RedirectIfUnauth(w http.ResponseWriter, r *http.Request) *sessions.Session {
	sess, err := s.SessStore.Get(r, model.SessionName)
	if err != nil && !sess.IsNew {
		s.RedirectError(w, r, model.ErrorSessionError)
		return nil
	} // Else implies it is a new session.
	return sess
}

// IsAuth checks whether user information is available in the session.
func (s *Server) IsAuth(r *http.Request) bool {
	sess, err := s.SessStore.Get(r, model.SessionName)
	if err != nil && !sess.IsNew {
		return false
	}
	authed := true
	if sess.Values[model.SessName] == nil || sess.Values[model.SessName] == "" {
		authed = false
	}
	return authed
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
