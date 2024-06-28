package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login handles the user authentication process using the CAS protocol and redirecting to an OAuth2 server.
// The function processes the following OPTIONAL query parameters from the request:
//   - service: URL of the service where the user intends to be redirected after authentication.
//   - renew: Indicates if the client wants to force re-authentication, regardless of existing session.
//   - gateway: If true, the client will not be prompted for credentials if not already logged in.
//
// Returns:
//   - Depending on the outcome of the authentication process and provided parameters,
//     the function may redirect the user, send specific error messages, or render certain views.
func Login(c *gin.Context) {
	span, _ := utils.StartAPMSpan(c.Request.Context(), config.AppConfig.UseAPM, utils.GetFunctionName(), "")
	defer utils.EndAPMSpan(span)

	serviceURL := c.DefaultQuery(constants.COMMON_SERVICE_PARAM, "")
	renew := c.DefaultQuery(constants.COMMON_RENEW_PARAM, "false")
	gateway := c.DefaultQuery(constants.COMMON_GATEWAY_PARAM, "false")

	utils.SetAPMLabel(span, constants.COMMON_SERVICE_PARAM, serviceURL)
	utils.SetAPMLabel(span, constants.COMMON_RENEW_PARAM, renew)
	utils.SetAPMLabel(span, constants.COMMON_GATEWAY_PARAM, gateway)

	isLoggedIn, username := isLoggedIn(c, config.AppConfig.TGTName)
	utils.SetAPMLabel(span, "isLoggedIn", isLoggedIn)

	if utils.IsTrue(gateway) && !utils.IsTrue(renew) {
		if isLoggedIn {
			redirectToService(c, serviceURL, username, "", false)
			return
		} else {
			c.Redirect(http.StatusSeeOther, serviceURL)
			return
		}
	}

	if isLoggedIn && !utils.IsTrue(renew) {
		redirectToService(c, serviceURL, username, "", false)
		return
	}

	if config.AppConfig.AuthMethod == constants.OAUTH_METHOD {
		setCookie(c, constants.SERVICE_URL_COOKIE, serviceURL, config.AppConfig.Domain, 3600)
		config.AuthProvider.RedirectAuth(c)
		return
	}
}

func isLoggedIn(c *gin.Context, tgtName string) (bool, string) {
	tgtCookie, err := c.Cookie(tgtName)
	if err != nil {
		return false, ""
	}

	return utils.ValidateTGT(tgtCookie)
}
