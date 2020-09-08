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
}

func routePages(s *Server, useMiddleware bool) {
	pages[model.PathDashboard] = &webpage{
		Title:       "Dashboard",
		Description: model.PageDescription,
		Path:        model.PathDashboard,
		Fn:          s.pageDashboard(),
	}
	pages[model.PathError] = &webpage{
		Title:       "Error",
		Description: model.PageDescription,
		Path:        model.PathError,
		Fn:          s.pageError(),
	}
	pages[model.PathIndex] = &webpage{
		Title:       "Home",
		Description: model.PageDescription,
		Path:        model.PathIndex,
		Fn:          s.pageIndex(),
	}
	pages[model.PathInitAdmin] = &webpage{
		Title:       "Initialization",
		Description: model.PageDescription,
		Path:        model.PathInitAdmin,
		Fn:          s.pageInitAdmin(),
	}
	pages[model.PathOpenIDCB] = &webpage{
		Title: "Oauth callback",
		Path:  model.PathOpenIDCB,
		Fn:    s.pageOpenIDCB(),
	}
	pages[model.PathVerifyClassroom] = &webpage{
		Title: "Test reading of classrooms",
		Path:  model.PathVerifyClassroom,
		Fn:    s.pageVerifyClassroom(),
	}
	pages[model.PathVerifySpreadsheet] = &webpage{
		Title: "Test reading of spreadsheet",
		Path:  model.PathVerifySpreadsheet,
		Fn:    s.pageVerifySpreadsheet(),
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
