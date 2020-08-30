package server

import (
	"net/http"
	"server-googleapi/google"
	"server-googleapi/model"
	"server-googleapi/tpl"
)

func (s *Server) pageError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "error"
		v := r.URL.Query()
		d := make(map[string]interface{})

		if v.Get("msg") != "" {
			d["msg"] = v.Get("msg")
		} else {
			d["msg"] = "Unknown server error"
		}

		tpl.Render(w, name, d)
	}
}

func (s *Server) pageIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "index"
		d := make(map[string]interface{})
		d["HasReadyCredentials"] = s.Ready
		cfg, err := google.MakeConfig(s.config.ReadGoogleCredentials(), nil)
		if err == nil {
			http.Redirect(w, r, "/error?msg="+model.ErrorOauthConstruction, http.StatusTemporaryRedirect)
			// If StatusInternalServerError, shows an ugly "Internal Server Error" page.
			return
		}
		d["UrlUserSignin"] = google.MakeLinkOnline(cfg, "csrf")
		tpl.Render(w, name, d)
	}
}
