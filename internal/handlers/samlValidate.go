package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SamlValidate "disguised" as a SAML endpoint.
// This function imitates the behavior of a SAML validation endpoint
// but internally uses CAS to OAuth2 validation mechanisms.
// Expected input and output formats are aligned with standard SAML responses.
// Parameters from query string:
//   - TARGET: The URL of the service requesting authentication.
//
// Parameters from body:
//   - SAMLRequest: The SAML request to validate in XML
//
// Returns:
//   - An XML response that either confirms the validity of the service ticket
//     or provides an error message indicating the reason for validation failure..
func SamlValidate(c *gin.Context) {
	serviceUrl := c.DefaultQuery(constants.SAML_TARGET_PARAM, "")
	var samlRequest SAMLRequest
	if err := c.ShouldBindXML(&samlRequest); err != nil {
		samlResponseError(c, constants.SAML_ERRMSG_INVALID_REQUEST)
		return
	}

	isValid, username, _, isOk := validationSAML(c, samlRequest.Body.Request.AssertionArtifact, serviceUrl)
	if !isOk {
		samlResponseError(c, constants.SAML_ERRMSG_VALIDATION)
		return
	}

	if !isValid {
		samlResponseError(c, constants.SAML_ERRMSG_INVALID_TICKET)
		return
	}

	samlResponseSuccess(c, serviceUrl, username)
}

func validationSAML(c *gin.Context, serviceTicket string, serviceURL string) (bool, string, bool, bool) {
	span, _ := utils.StartAPMSpan(c.Request.Context(), config.AppConfig.UseAPM, utils.GetFunctionName(), "")
	defer utils.EndAPMSpan(span)

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

func samlResponseError(c *gin.Context, message string) {
	t1 := time.Now().UTC().Truncate(time.Millisecond)
	response := ResponseEnvelope{
		Body: ResponseBody{
			Response: SAMLResponse{
				ResponseID:   utils.RandomString(16),
				MajorVersion: 1,
				MinorVersion: 1,
				IssueInstant: t1,
				Status: Status{
					StatusCode: StatusCode{
						Value: constants.SAML_STATUSCODE_ERROR,
					},
				},
			},
		},
	}

	c.XML(http.StatusForbidden, response)
}

func samlResponseSuccess(c *gin.Context, serviceUrl, username string) {
	t1 := time.Now().UTC().Truncate(time.Millisecond)

	response := ResponseEnvelope{
		XMLNS: constants.XML_SOAP_NAMESPACE,
		Body: ResponseBody{
			Response: SAMLResponse{
				XMLNS:        fmt.Sprintf("%s:%s", constants.XML_SAML_NAMESPACE, "protocol"),
				ResponseID:   fmt.Sprintf("_%s", utils.RandomString(17)),
				Recipient:    serviceUrl,
				MajorVersion: 1,
				MinorVersion: 1,
				IssueInstant: t1,
				Status: Status{
					StatusCode: StatusCode{
						Value: constants.SAML_STATUSCODE_SUCCESS,
					},
				},
				Assertion: Assertion{
					XMLNS:        fmt.Sprintf("%s:%s", constants.XML_SAML_NAMESPACE, "assertion"),
					AssertionID:  fmt.Sprintf("_%s", utils.RandomString(16)),
					Issuer:       constants.SAML_ISSUER,
					IssueInstant: t1,
					MajorVersion: 1,
					MinorVersion: 1,
					Conditions: Conditions{
						NotBefore:    t1,
						NotOnOrAfter: t1.Add(time.Minute * 1),
						AudienceRestrictionCondition: AudienceRestrictionCondition{
							Audience: serviceUrl,
						},
					},
					AuthenticationStatement: AuthenticationStatement{
						AuthenticationMethod:  fmt.Sprintf("%s:%s", constants.XML_SAML_NAMESPACE, "am:unspecified"),
						AuthenticationInstant: t1,
						Subject: Subject{
							NameIdentifier: username,
							SubjectConfirmation: SubjectConfirmation{
								ConfirmationMethod: fmt.Sprintf("%s:%s", constants.XML_SAML_NAMESPACE, "cm:artifact"),
							},
						},
					},
				},
			},
		},
	}

	c.XML(http.StatusOK, response)
}
