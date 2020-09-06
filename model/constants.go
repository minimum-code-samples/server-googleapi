package model

const (
	// ErrorConfigConstruction describes the error when creating the Oauth config object.
	ErrorConfigConstruction = "An error occurred while creating Oauth config. Please contact the system administrator."
	// ErrorJWTDecode describes the error when parsing a JWT from the Oauth token fails.
	ErrorJWTDecode = "An error occurred while parsing user information. Please contact the system administrator."
	// ErrorOauthConstruction is the error message when creating the Oauth object fails.
	ErrorOauthConstruction = "An error occurred while formulating Oauth authentication. Please contact the system administrator."
	// ErrorOauthExchange is the error message when exchange for the Oauth token fails.
	ErrorOauthExchange = "An error occurred while requesting access token."
	// ErrorSessionError is the generic error related to sessions.
	ErrorSessionError = "An occurred for the session. Please contact the system administrator."
	// ErrorSessionSave logs when saving data to the session results in an error.
	ErrorSessionSave = "Unable to save to session: %s"
	// ErrorTokenSave is the error message when the saving of the token fails.
	ErrorTokenSave = "Failed to complete server initialization."
	// PathDashboard is the user start page.
	PathDashboard = "/dashboard"
	// PageDescription contains the value for the description meta tag.
	PageDescription = "Sample Server"
	// PathError is the general error page.
	PathError = "/error"
	// PathIndex is the home page.
	PathIndex = "/"
	// PathInitAdmin is the page to start initialization for the server.
	PathInitAdmin = "/init-admin"
	// PathOpenIDCB is the endpoint for the OpenID redirect callback.
	//
	// If this needs to be changed, the URI in the Client ID setting (https://console.cloud.google.com/apis/credentials) needs to be updated correspondingly.
	PathOpenIDCB = "/openidcb"
	// SessName is the session variable for the user's name.
	SessName = "name"
	// SessionName is the name of the session.
	SessionName = "minserver"
)
