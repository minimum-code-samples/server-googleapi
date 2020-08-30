package server

import (
	"net/http"
	"server-googleapi/tpl"
)

func (s *Server) pageIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "index"
		d := make(map[string]interface{})
		tpl.Render(w, name, d)
	}
}
