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

func (s *Server) pageError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "error"
		v := r.URL.Query()
		d := make(map[string]interface{})

		if v.Get(model.QueryMsg) != "" {
			d[model.QueryMsg] = v.Get(model.QueryMsg)
		} else {
			d[model.QueryMsg] = "Unknown server error"
		}

		tpl.Render(w, name, d)
	}
}

func (s *Server) pageIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "index"
		d := make(map[string]interface{})
		d["HasReadyCredentials"] = s.TokenAdmin != nil
		if s.IsAuth(r) {
			d["Authed"] = true
		} else {
			var (
				cfg *oauth2.Config
				err error
			)
			cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithSheets())
			if err != nil {
				lg.Error(lg.CriticalOauthConfig, err)
				s.RedirectError(w, r, model.ErrorOauthConstruction)
				return
			}
			if s.TokenAdmin == nil {
				// Create request for a new refresh token.
				d["UrlUserSignin"] = google.MakeLinkOffline(cfg, "csrf")
			} else {
				d["UrlUserSignin"] = google.MakeLinkOnline(cfg, "csrf")
			}
		}
		tpl.Render(w, name, d)
	}
}

// PageOpenIDCB is the handler for "/openidcb" to handle the OpenID callback.
func (s *Server) pageOpenIDCB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		queries := r.URL.Query()
		code := queries.Get("code")
		_ = queries.Get("state")
		var (
			cfg *oauth2.Config
			err error
		)
		cfg, err = google.MakeConfig(s.config.ReadGoogleCredentials(), google.ScopesWithSheets())
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

		if token.RefreshToken != "" {
			// i.e. As an administrator.
			http.Redirect(w, r, model.PathVerifySpreadsheetAdmin, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, model.PathVerifySpreadsheet, http.StatusTemporaryRedirect)
		return
	}
}

func (s *Server) pageVerifySpreadsheet(asAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := "verify-spreadsheet"
		d := make(map[string]interface{})

		creds := s.config.ReadGoogleCredentials()
		var tok *oauth2.Token
		if asAdmin {
			tok = s.TokenAdmin
		} else {
			tok = s.ReadToken(w, r)
		}

		// Get the spreadsheet ID.
		q := r.URL.Query()
		sid := q.Get(model.QuerySpreadsheetID)
		sn := q.Get(model.QuerySheetName)
		if sid == "" || sn == "" {
			d["Error"] = fmt.Sprintf("Append query parameters '%s' and '%s'", model.QuerySpreadsheetID, model.QuerySheetName)
		} else {
			values, err := google.FetchSpreadsheetValues(ctx, sid, sn, creds, tok)
			if err != nil {
				d["Error"] = err.Error()
			} else {
				d["Count"] = len(values)
				fmt.Printf("1st dimension of results is %d long.", len(values))
			}
		}
		tpl.Render(w, name, d)
	}
}
