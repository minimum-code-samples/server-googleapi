package lg

const (
	// CriticalOauthConfig logs when the `MakeConfig` function returns an error. It is critical because the system will not be able to perform its expected function but does not crash it.
	CriticalOauthConfig = "Unable to make Oauth config: %s"
	// CriticalOauthDecode logs when decoding the JWT from the token fails.
	CriticalOauthDecode = "Unable to decode user info from Oauth token: %s"
	// CriticalOauthExchange logs when the Oauth exchange process fails.
	CriticalOauthExchange = "Unable to complete the Oauth token exchange process: %s"
	// CriticalTokenSave logs when saving of the token fails.
	CriticalTokenSave = "Unable to save token: %s"
	// FatalConfigParse logs when the config file cannot be parsed.
	FatalConfigParse = "Unable to parse config file:\n%s\n"
	// FatalGACEmpty logs when the configuration property is empty and not specified via a command-line flag.
	FatalGACEmpty = "'google_application_credentials' not specified."
	// FatalGATEmpty logs when the configuration property is empty.
	FatalGATEmpty = "'google_admin_token' not specified."
	// FatalGACParse logs when the GAC
	FatalGACParse = "Unable to parse GOOGLE_APPLICATIONS_FILE:\n%s\n"
	// FatalSessionKeyEmpty logs when the authentication key to the cookie is not set.
	FatalSessionKeyEmpty = "Authentication key for cookie must be present."
	// FatalTokenFileCorrupt logs when the file containing the token is not valid.
	FatalTokenFileCorrupt = "File containing admin token is valid or corrupted. Remove or replace with a valid one."
	// ServerInit logs when the server is starting.
	ServerInit = "Initializing server on %s:%s"
	// ServerInitSecure logs when the server is starting on two ports.
	ServerInitSecure = "Initializing server on %s:%s and %[1]s:%[3]s"
	// ServerStarted logs when the server has started.
	ServerStarted = "Server started."
)
