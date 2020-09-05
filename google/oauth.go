package google

import (
	"context"
	"encoding/json"
	"os"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/sheets/v4"
)

const (
	openIDIssuer = "https://accounts.google.com"
)

// UserDetail holds the detailed information from the UserInfo provider.
type UserDetail struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
	HD            string `json:"hd"`
}

// DeriveUserInfo retrieves the user information using the supplied token.
//
// This is part of the OpenID Connect "server flow".
func DeriveUserInfo(ctx context.Context, token *oauth2.Token) (userInfo *oidc.UserInfo, userDetail *UserDetail, err error) {
	userDetail = new(UserDetail)
	provider, err := oidc.NewProvider(ctx, openIDIssuer)
	if err != nil {
		return
	}
	userInfo, err = provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return
	}

	if err = userInfo.Claims(userDetail); err != nil {
		return
	}
	return
}

// SaveTokenAsFile saves the token into the filesystem.
func SaveTokenAsFile(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

// ScopesWithClassroom creates the scopes with the necessary Classroom scopes added.
func ScopesWithClassroom() []string {
	return []string{"profile", "email", classroom.ClassroomCoursesReadonlyScope}
}

// ScopesWithSheets creates the scopes along with readonly access to Google Sheets.
func ScopesWithSheets() []string {
	return []string{"profile", "email", sheets.SpreadsheetsReadonlyScope}
}

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
