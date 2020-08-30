package lg

const (
	// FatalConfigParse logs when the config file cannot be parsed.
	FatalConfigParse = "Unable to parse config file:\n%s\n"
	// FatalGACParse logs when the GAC
	FatalGACParse = "Unable to parse GOOGLE_APPLICATIONS_FILE:\n%s\n"
	// ServerInit logs when the server is starting.
	ServerInit = "Initializing server on %s:%s"
	// ServerInitSecure logs when the server is starting on two ports.
	ServerInitSecure = "Initializing server on %s:%s and %[1]s:%[3]s"
	// ServerStarted logs when the server has started.
	ServerStarted = "Server started."
)
