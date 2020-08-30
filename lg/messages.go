package lg

const (
	// FatalConfigParse logs when the config file cannot be parsed.
	FatalConfigParse = "Unable to parse config file:\n%s\n"
	// FatalGACEmpty logs when the configuration property is empty and not specified via a command-line flag.
	FatalGACEmpty = "'google_application_credentials' not specified."
	// FatalGACParse logs when the GAC
	FatalGACParse = "Unable to parse GOOGLE_APPLICATIONS_FILE:\n%s\n"
	// ServerInit logs when the server is starting.
	ServerInit = "Initializing server on %s:%s"
	// ServerInitSecure logs when the server is starting on two ports.
	ServerInitSecure = "Initializing server on %s:%s and %[1]s:%[3]s"
	// ServerStarted logs when the server has started.
	ServerStarted = "Server started."
)
