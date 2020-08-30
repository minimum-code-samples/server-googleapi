package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// MakeConfig creates the configuration object with the necessary client information to perform a token exchange.
func MakeConfig(file []byte, scopes []string) (*oauth2.Config, error) {
	return google.ConfigFromJSON(file, scopes...)
}

// MakeLinkOffline creates an Oauth link for "offline" access i.e. request for refresh token.
func MakeLinkOffline(cfg *oauth2.Config, csrfToken string) string {
	igs := oauth2.SetAuthURLParam("include_granted_scopes", "true")
	at := oauth2.SetAuthURLParam("access_type", "offline")
	// Need consent in case user did not grant offline access the first time.
	p := oauth2.SetAuthURLParam("prompt", "consent")
	return cfg.AuthCodeURL(csrfToken, igs, at, p)
}

// MakeLinkOnline creates an Oauth link.
func MakeLinkOnline(cfg *oauth2.Config, csrfToken string) string {
	// hd := oauth2.SetAuthURLParam("hd", "*")
	igs := oauth2.SetAuthURLParam("include_granted_scopes", "true")
	return cfg.AuthCodeURL(csrfToken, igs)
}
