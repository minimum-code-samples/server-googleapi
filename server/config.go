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
	// The interface that the server is listening on.
	Interface string `yaml:"interface"`
	// The folder where the logs are placed in.
	LogFolder string `yaml:"log_folder"`
	// LogLevel is the lowest level of logging that is saved to the log file.
	//
	// Regardless of the level, all log statements are sent to STDOUT.
	LogLevel string `yaml:"log_level"`
	// The Google Oauth client credentials.
	GoogleApplicationCredentials string `yaml:"google_application_credentials"`
	// The port that the server will run on.
	Port string `yaml:"port"`
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
		log.Fatalf(lg.FatalGACEmpty)
	}
	file, err := ioutil.ReadFile(c.GoogleApplicationCredentials)
	if err != nil {
		log.Fatalf(lg.FatalGACParse, err)
	}
	return file
}
