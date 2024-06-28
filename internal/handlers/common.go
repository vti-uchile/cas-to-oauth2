package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func redirectToService(c *gin.Context, serviceURL, username, tgt string, isDirect bool) {
	if serviceURL == "" {
		c.HTML(http.StatusOK, constants.LOGIN_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_OK})
		return
	}

	isAllowed := checkAllowedDomains(serviceURL)
	if !isAllowed {
		c.HTML(http.StatusForbidden, constants.UNAUTHORIZED_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.COMMON_ERRMSG_INVALID_SERVICE})
		return
	}

	var err error
	if tgt == "" {
		tgt, err = c.Cookie(config.AppConfig.TGTName)
		if err != nil {
			c.HTML(http.StatusBadRequest, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.COMMON_ERRMSG_MISSING})
			return
		}
	}

	serviceTicket := utils.GenerateServiceTicket(serviceURL, username, tgt, isDirect)
	redirectURL := fmt.Sprintf("%s?ticket=%s", serviceURL, serviceTicket)
	c.Redirect(http.StatusFound, redirectURL)
}

func checkAllowedDomains(serviceURL string) bool {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return false
	}

	host := u.Hostname()
	for _, domain := range config.AppConfig.AllowedDomains {
		if strings.HasSuffix(host, domain) {
			return true
		}
	}

	return false
}

func setCookie(c *gin.Context, tgtName, tgtValue, domain string, duration int) {
	c.SetCookie(tgtName, tgtValue, duration, "/", domain, config.AppConfig.TGTSecure, config.AppConfig.TGTHttpOnly)
}

func unsetCookie(c *gin.Context, tgtName, domain string) {
	c.SetCookie(tgtName, "", -1, "/", domain, config.AppConfig.TGTSecure, config.AppConfig.TGTHttpOnly)
}
