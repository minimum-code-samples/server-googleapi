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
	// ErrorSessionToken logs when the token cannot be retrieved from the session cookie.
	ErrorSessionToken = "Corrupted/missing token in session."
	// ErrorSessionUnauth is the error message for a user who is not signed in.
	ErrorSessionUnauth = "User not signed in."
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
	// PathOpenIDCB is the endpoint for the OpenID redirect callback.
	//
	// If this needs to be changed, the URI in the Client ID setting (https://console.cloud.google.com/apis/credentials) needs to be updated correspondingly.
	PathOpenIDCB = "/openidcb"
	// PathVerifySpreadsheet is a temporary endpoint.
	// TODO Change this.
	PathVerifySpreadsheet = "/verify-spreadsheet"
	// PathVerifySpreadsheetAdmin is a temporary endpoint.
	// TODO Change this.
	PathVerifySpreadsheetAdmin = "/verify-spreadsheet-admin"
	// QueryMsg is the query parameter specifying the error message.
	QueryMsg = "msg"
	// QuerySheetName is the query parameter specifying the name of the sheet to retrieve.
	QuerySheetName = "sheet-name"
	// QuerySpreadsheetID is the query parameter specifying the ID of the Google Sheets spreadsheet.
	QuerySpreadsheetID = "spreadsheet-id"
	// SessName is the session variable for the user's name.
	SessName = "name"
	// SessToken is the session variable for the user's access token.
	SessToken = "token"
	// SessionName is the name of the session.
	SessionName = "minserver"
)
