package server

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"server-googleapi/lg"
)

// Config stores the configuration for the server.
type Config struct {
	// The key to create CSRF tokens.
	CSRFKey string `yaml:"csrf_key"`
	// The token for admin access.
	GoogleAdminToken string `yaml:"google_admin_token"`
	// The path to the file containing Google Oauth client credentials.
	GoogleApplicationCredentials string `yaml:"google_application_credentials"`
	// The interface that the server is listening on.
	Interface string `yaml:"interface"`
	// The folder where the logs are placed in.
	LogFolder string `yaml:"log_folder"`
	// LogLevel is the lowest level of logging that is saved to the log file.
	//
	// Regardless of the level, all log statements are sent to STDOUT.
	LogLevel string `yaml:"log_level"`
	// The port that the server will run on.
	Port string `yaml:"port"`
	// The authentication key. Required.
	SessionAuthKey string `yaml:"session_auth_key"`
	// The number of minutes that the session is valid for.
	SessionDuration int `yaml:"session_duration"`
	// The encryption key to obfuscate the cookie values. Optional.
	SessionEncKey string `yaml:"session_enc_key"`
	// Whether to serve secure sessions only.
	SessionSecure bool `yaml:"session_secure_only"`
}

// MakeCSRFToken creates a CSRF token based on the `input` and an optional `suffix`.
//
// The token is based on the `CSRFKey` specified in the Config.
//
// Returns a Base64 encoded string.
func (c *Config) MakeCSRFToken(input, suffix string) string {
	src := c.CSRFKey + ":" + input
	if suffix != "" {
		src = src + ":" + suffix
	}
	keyHash := sha256.Sum256([]byte(src))
	return base64.StdEncoding.EncodeToString(keyHash[:])
}

// ReadGoogleCredentials reads the credentials file.
func (c *Config) ReadGoogleCredentials() []byte {
	if c.GoogleApplicationCredentials == "" {
		log.Fatal(lg.FatalGACEmpty)
	}
	file, err := ioutil.ReadFile(c.GoogleApplicationCredentials)
	if err != nil {
		log.Fatalf(lg.FatalGACParse, err)
	}
	return file
}
