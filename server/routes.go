package server

import (
	"net/http"
	"server-googleapi/model"
)

var (
	pages = make(map[string]*webpage)
)

type webpage struct {
	Title       string
	Description string
	Path        string
	Fn          http.HandlerFunc
	Template    string
}

func routePages(s *Server, useMiddleware bool) {
	pages[model.PathIndex] = &webpage{
		Title:       "Home",
		Description: model.PageDescription,
		Path:        model.PathIndex,
		Fn:          s.pageIndex(),
		Template:    "templates/web/index.html",
	}

	for _, p := range pages {
		fn := p.Fn
		// if useMiddleware {
		// 	fn := middleware.Chain(p.Fn, p.Middleware...)
		// }
		s.Router.HandleFunc(p.Path, fn)
	}
}

func routeStatic(s *Server) {
	// TODO Add an error handling middleware.
	fileServer := http.FileServer(http.Dir("./static"))
	// TODO Add caching middleware.
	s.Router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
}
