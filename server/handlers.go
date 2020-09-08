package server

import (
	"fmt"
	"net/http"
	"server-googleapi/google"
	"server-googleapi/lg"
	"server-googleapi/model"
	"server-googleapi/tpl"

	"golang.org/x/oauth2"
)

func (s *Server) pageDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "dashboard"
		d := make(map[string]interface{})
		tpl.Render(w, name, d)
	}
}

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
		if s.IsAuth(r) {
			d["Authed"] = true
		} else {
			d["HasReadyCredentials"] = s.TokenAdmin != nil
			var (
				cfg *oauth2.Config
				err error
			)
			if s.TokenAdmin == nil {
				// Create request for a new refresh token.
				cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithSheets())
			} else {
				cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithClassroom())
			}
			if err != nil {
				lg.Error(lg.CriticalOauthConfig, err)
				s.RedirectError(w, r, model.ErrorOauthConstruction)
				return
			}
			if s.TokenAdmin == nil {
				d["UrlUserSignin"] = google.MakeLinkOffline(cfg, "csrf")
			} else {
				d["UrlUserSignin"] = google.MakeLinkOnline(cfg, "csrf")
			}
		}
		tpl.Render(w, name, d)
	}
}

func (s *Server) pageInitAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "init-admin"
		d := make(map[string]interface{})
		tpl.Render(w, name, d)
	}
}

// PageOpenIDCB is the handler for "/openidcb" to handle the OpenID callback.
func (s *Server) pageOpenIDCB() http.HandlerFunc {
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
			s.RedirectError(w, r, model.ErrorConfigConstruction)
			return
		}
		token, err := cfg.Exchange(ctx, code)
		if err != nil {
			lg.Error(lg.CriticalOauthExchange, err)
			s.RedirectError(w, r, model.ErrorOauthExchange)
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
			s.RedirectError(w, r, model.ErrorJWTDecode)
			return
		}
		lg.Debug("Name:\n  %v\n", userDetail.Name)
		lg.Debug("Email:\n  %v\n", userDetail.Email)

		// Store the refresh token for the admin account.
		if token.RefreshToken != "" {
			if err = google.SaveTokenAsFile(s.config.GoogleAdminToken, token); err != nil {
				lg.Error(lg.CriticalTokenSave, err)
				s.RedirectError(w, r, model.ErrorTokenSave)
			}
		}

		// Store in session.
		if token.RefreshToken == "" {
			// Store access token in session cookie as well.
			if !s.AuthToken(userDetail, token, w, r) {
				return // Redirection already done in Auth().
			}
		} else {
			// Don't store token if it contains refresh token.
			if !s.Auth(userDetail, w, r) {
				return // Redirection already done in Auth().
			}
		}

		// TODO Check state to determine redirect destination.
		if state != "" {
			http.Redirect(w, r, model.PathInitAdmin, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, model.PathDashboard, http.StatusTemporaryRedirect)
		return
	}
}

func (s *Server) pageVerifyClassroom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		d := make(map[string]interface{})
		token := s.ReadToken(w, r)
		if token == nil {
			return // Redirection already done in ReadToken.
		}
		creds := s.config.ReadGoogleCredentials()
		if kourses, err := google.FetchCourses(ctx, creds, token); err != nil {
			d["Error"] = err.Error()
		} else {
			d["Titles"] = google.ReadClassroomNames(kourses)
		}
		tpl.Render(w, "verify-spreadsheet", d)
	}
}

func (s *Server) pageVerifySpreadsheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := "verify-spreadsheet"
		d := make(map[string]interface{})
		creds := s.config.ReadGoogleCredentials()
		spreadsheet := "1Mt2AQLBUfZ9ZAmCBP-6X3aFx3RJ5rUmor02iHVI64sU"
		fmt.Printf("- TokenAdmin: %v\n", s.TokenAdmin)
		if sheeds, err := google.FetchSpreadsheetSheets(ctx, spreadsheet, creds, s.TokenAdmin); err != nil {
			d["Error"] = err.Error()
		} else {
			d["Titles"] = google.ReadSheetsTitles(sheeds)
		}

		tpl.Render(w, name, d)
	}
}
