package constants

const (
	// Endpoints
	ENDPOINT_ROOT             = "/"
	ENDPOINT_LOGIN            = "/login"
	ENDPOINT_OAUTH2           = "/oauth2/callback"
	ENDPOINT_SERVICE_VALIDATE = "/serviceValidate"
	ENDPOINT_PROXY_VALIDATE   = "/proxyValidate"
	ENDPOINT_SAML_VALIDATE    = "/samlValidate"
	ENDPOINT_VALIDATE         = "/validate"
	ENDPOINT_PROXY            = "/proxy"
	ENDPOINT_LOGOUT           = "/logout"
	ENDPOINT_HEALTHCHECK      = "/healthcheck"

	// Template variables
	TEMPLATE_MESSAGE = "message"

	// Cookies
	SERVICE_URL_COOKIE = "serviceURL"

	// Main
	MAIN_ERRMSG = "Error starting server"

	// Common
	COMMON_SERVICE_PARAM          = "service"
	COMMON_RENEW_PARAM            = "renew"
	COMMON_GATEWAY_PARAM          = "gateway"
	COMMON_ERRMSG_MISSING         = "Ticket Granting Ticket is missing"
	COMMON_ERRMSG_URL_PARSE       = "Error parsing URL"
	COMMON_ERRMSG_INVALID_SERVICE = "Service access is not allowed"

	// OAuth2Callback
	OAUTH_METHOD               = "oauth2"
	OAUTH_CODE_PARAM           = "code"
	OAUTH_ERRMSG_UNAUTHORIZED  = "Authorization code missing"
	OAUTH_ERRMSG_INVALID_TOKEN = "Invalid token"
	OAUTH_ERRMSG_EXCHANGE      = "Error exchanging code for token"
	OAUTH_ERRMSG_SUB           = "Error getting subject from token"
	OAUTH_ERRMSG_OK            = "TGT successfully generated"
	OAUTH_ERRMSG_SPAN          = "Return from OAuth2 provider"

	// ServiceValidate
	VALIDATE_TICKET_PARAM           = "ticket"
	VALIDATE_INVALID_REQUEST        = "INVALID_REQUEST"
	VALIDATE_INVALID_TICKET         = "INVALID_TICKET"
	VALIDATE_XML_RESPONSE           = "Error generating XML response"
	VALIDATE_ERRMSG_INVALID_REQUEST = "Service Ticket or Service URL is missing"
	VALIDATE_ERRMSG_INVALID_TICKET  = "Invalid Service Ticket"
	VALIDATE_IS_VALID               = "IsSTValid"
	VALIDATE_IS_DIRECT              = "IsSTDirect"

	// Logout
	LOGOUT_REDIRECT_PARAM    = "url"
	LOGOUT_ERRMSG_MISSING    = "TGT Cookie is missing"
	LOGOUT_ERRMSG_DELETE_TGT = "Error deleting TGT"
	LOGOUT_OK                = "TGT successfully deleted"

	// Utils
	UTILS_ID_TOKEN               = "id_token"
	UTILS_CLAIM                  = "sub"
	UTILS_ERRMSG_MISSING         = "No id_token field in oauth2 token"
	UTILS_ERRMSG_CLAIM_PARSE     = "Claims could not be parsed"
	UTILS_ERRMSG_CLAIM_NOT_EXIST = "Claim does not exist"

	// Templates
	UNAUTHORIZED_HTML = "unauthorized.html"
	ERROR_HTML        = "error.html"
	LOGIN_HTML        = "login.html"
	LOGOUT_HTML       = "logout.html"

	// CAS XML Namespaces
	XML_CAS_NAMESPACE = "http://www.yale.edu/tp/cas"

	// Database Collections
	DB_COLLECTION_SERVICE_TICKETS = "serviceTickets"
	DB_COLLECTION_TGT             = "ticketGrantingTickets"

	// SAML Validate
	SAML_TARGET_PARAM           = "TARGET"
	SAML_ERRMSG_INVALID_REQUEST = "Invalid SAML Request"
	SAML_ERRMSG_VALIDATION      = "Error in validation process"
	SAML_ERRMSG_INVALID_TICKET  = "Invalid SAML Ticket or Service"
	SAML_ISSUER                 = "cas-to-oauth2"
	SAML_STATUSCODE_SUCCESS     = "saml1p:Success"
	SAML_STATUSCODE_ERROR       = "saml1p:RequestDenied"
	XML_SOAP_NAMESPACE          = "http://schemas.xmlsoap.org/soap/envelope/"
	XML_SAML_NAMESPACE          = "urn:oasis:names:tc:SAML:1.0"
)
