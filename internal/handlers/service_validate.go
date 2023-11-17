package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceValidate validates the service ticket provided by the client.
// It processes the request based on query string parameters and returns an XML response.
// Parameters from query string:
//   - ticket: The service ticket issued by the CAS server.
//   - service: The URL of the service requesting authentication.
//   - renew(optional): Indicates whether to force re-authentication, ignoring single sign-on sessions
//
// Returns:
//   - An XML response that either confirms the validity of the service ticket
//     or provides an error message indicating the reason for validation failure.
func ServiceValidate(c *gin.Context) {
	var response CASResponse
	response.XMLNS = constants.XML_CAS_NAMESPACE
	formatted := true

	isValid, username, _, isOk := commonValidation(c)
	if !isOk {
		response.Failure = &AuthenticationFailure{Code: constants.VALIDATE_INVALID_REQUEST, Description: constants.VALIDATE_ERRMSG_INVALID_REQUEST}
		xmlResponse(c, http.StatusUnauthorized, response, formatted)
		return
	}

	if !isValid {
		response.Failure = &AuthenticationFailure{Code: constants.VALIDATE_INVALID_TICKET, Description: constants.VALIDATE_ERRMSG_INVALID_TICKET}
		xmlResponse(c, http.StatusUnauthorized, response, formatted)
		return
	}

	response.Success = &AuthenticationSuccess{User: username}
	xmlResponse(c, http.StatusOK, response, formatted)
}

// Validate validates the service ticket provided by the client.
// It processes the request based on query string parameters and returns a plain text response.
// Parameters from query string:
//   - ticket: The service ticket issued by the CAS server.
//   - service: The URL of the service requesting authentication.
//   - renew(optional): Indicates whether to force re-authentication, ignoring single sign-on sessions
//
// Returns:
//   - A plain text response that either confirms the validity of the service ticket
//     or provides an error message indicating the reason for validation failure.
func Validate(c *gin.Context) {
	isValid, username, _, isOk := commonValidation(c)
	if !isOk {
		c.String(http.StatusBadRequest, "no\n")
		return
	}

	if !isValid {
		c.String(http.StatusUnauthorized, "no\n")
		return
	}

	c.String(http.StatusOK, "yes\n%s\n", username)
}

func commonValidation(c *gin.Context) (bool, string, bool, bool) {
	span, _ := utils.StartAPMSpan(c.Request.Context(), config.AppConfig.UseAPM, utils.GetFunctionName(), "")
	defer utils.EndAPMSpan(span)

	serviceTicket := c.DefaultQuery(constants.VALIDATE_TICKET_PARAM, "")
	serviceURL := c.DefaultQuery(constants.COMMON_SERVICE_PARAM, "")
	renew := c.DefaultQuery(constants.COMMON_RENEW_PARAM, "false")

	utils.SetAPMLabel(span, constants.COMMON_SERVICE_PARAM, serviceURL)
	utils.SetAPMLabel(span, constants.COMMON_RENEW_PARAM, renew)

	if serviceTicket == "" || serviceURL == "" {
		return false, "", false, false
	}

	isValid, username, isDirect := utils.ValidateServiceTicket(serviceTicket, serviceURL)
	if !isValid || (utils.IsTrue(renew) && !isDirect) {
		return false, "", false, true
	}

	utils.SetAPMLabel(span, constants.VALIDATE_IS_VALID, isValid)
	utils.SetAPMLabel(span, constants.VALIDATE_IS_DIRECT, isDirect)

	return true, username, isDirect, true
}

func xmlResponse(c *gin.Context, code int, response interface{}, formatted bool) {
	if formatted {
		xmlData, err := xml.MarshalIndent(response, "", "    ")
		if err != nil {
			c.HTML(http.StatusInternalServerError, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.VALIDATE_XML_RESPONSE})
			return
		}
		c.Data(code, "application/xml; charset=utf-8", xmlData)
		return
	}

	c.XML(code, response)
}
