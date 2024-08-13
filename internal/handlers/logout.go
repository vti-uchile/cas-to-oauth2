package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout handles the user logout process, invalidating the current user session.
// It processes the request primarily based on cookies and query string parameters.
// Cookies:
//   - TGTName: A cookie containing the Ticket Granting Ticket, used for CAS authentication.
//
// Query string parameters (optional):
//   - url(optional): A URL to redirect the user to after successful logout.
//
// Returns:
//   - Depending on the outcome, the function may redirect the user to a specified URL
//     or render a confirmation message of successful logout.
func Logout(c *gin.Context) {
	span, _ := utils.StartAPMSpan(c.Request.Context(), config.AppConfig.UseAPM, utils.GetFunctionName(), "")
	defer utils.EndAPMSpan(span)

	utils.SetAPMLabel(span, constants.COMMON_SERVICE_PARAM, c.Request.RequestURI)

	var err error
	tgtCookie, err := c.Cookie(config.AppConfig.TGTName)
	if err != nil {
		c.HTML(http.StatusBadRequest, constants.UNAUTHORIZED_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.LOGOUT_ERRMSG_MISSING})
		return
	}

	err = utils.DeleteTGT(tgtCookie)
	if err != nil {
		c.HTML(http.StatusBadRequest, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.LOGOUT_ERRMSG_DELETE_TGT})
		return
	}

	unsetCookie(c, config.AppConfig.TGTName, config.AppConfig.Domain)
	unsetCookie(c, config.AppConfig.JSessionID, config.AppConfig.Domain)

	if config.AppConfig.AnotherCookie != "" {
		unsetCookie(c, config.AppConfig.AnotherCookie, config.AppConfig.Domain)
	}

	if url := c.Query(constants.LOGOUT_REDIRECT_PARAM); url != "" {
		c.Redirect(http.StatusFound, url)
		return
	}

	if url := c.Query(constants.COMMON_SERVICE_PARAM); url != "" {
		c.Redirect(http.StatusFound, url)
		return
	}

	loginURL := c.Request.URL.Scheme + "://" + c.Request.URL.Host + constants.ENDPOINT_LOGIN
	c.Redirect(http.StatusFound, loginURL)
}
