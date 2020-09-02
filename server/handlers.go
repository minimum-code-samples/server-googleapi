package server

import (
	"encoding/json"
	"net/http"
	"server-googleapi/google"
	"server-googleapi/lg"
	"server-googleapi/model"
	"server-googleapi/tpl"

	"golang.org/x/oauth2"
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
		var (
			cfg *oauth2.Config
			err error
		)
		if !s.Ready {
			cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithSheets())
		} else {
			cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithClassroom())
		}
		if err != nil {
			lg.Error(lg.CriticalOauthConfig, err)
			http.Redirect(w, r, "/error?msg="+model.ErrorOauthConstruction, http.StatusTemporaryRedirect)
			// If StatusInternalServerError, shows an ugly "Internal Server Error" page.
			return
		}
		d["UrlUserSignin"] = google.MakeLinkOnline(cfg, "csrf")
		tpl.Render(w, name, d)
	}
}

// PageOpenIDCB is the handler for "/openidcb" to handle the OpenID callback.
func (s *Server) PageOpenIDCB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		queries := r.URL.Query()
		code := queries.Get("code")
		state := queries.Get("state")
		var (
			cfg *oauth2.Config
			err error
		)
		// TODO Check state to determine the scope to exchange token.
		if state != "" {
			// Treat this as an admin account.
			cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithSheets())
		} else {
			// Handle as "normal" user.
			cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithClassroom())

		}
		if err != nil {
			lg.Error(lg.CriticalOauthConfig, err)
			http.Redirect(w, r, "/error?msg="+model.ErrorConfigConstruction, http.StatusTemporaryRedirect)
			return
		}
		token, err := cfg.Exchange(ctx, code)
		if err != nil {
			lg.Error(lg.CriticalOauthExchange, err)
			http.Redirect(w, r, "/error?msg="+model.ErrorOauthExchange, http.StatusTemporaryRedirect)
			return
		}
		lg.Debug("Access token:\n  %v\n", token.AccessToken)
		lg.Debug("Token type:\n  %v\n", token.TokenType)
		lg.Debug("Refresh token:\n  %v\n", token.RefreshToken)
		lg.Debug("Expiry:\n  %v\n", token.Expiry)
		idToken := token.Extra("id_token").(string)
		lg.Debug("idToken:\n  %v\n", idToken)
		// Decode the ID token to get the user info.
		_, userDetail, err := google.DeriveUserInfo(ctx, token)
		if err != nil {
			lg.Error(lg.CriticalOauthDecode, err)
			http.Redirect(w, r, "/error?msg="+model.ErrorJWTDecode, http.StatusTemporaryRedirect)
			return
		}
		lg.Debug("Name:\n  %v\n", userDetail.Name)
		lg.Debug("Email:\n  %v\n", userDetail.Email)

		// Store the refresh token for the admin account.
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			lg.Error(lg.CriticalTokenMarshal, err)
		}
		lg.Debug("AccessToken: %s", tokenJSON)
	}
}
