package handlers

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuth2Callback handles the callback from the OAuth2 provider after user authentication.
// It processes the request based on query string parameters and cookies.
// Parameters from query string:
//   - code: The authorization code returned by the OAuth2 provider.
//
// Cookies:
//   - serviceUrl: A cookie containing the service url, used for redirection after successful authentication.
//
// Returns:
//   - Depending on the outcome of the OAuth2 token exchange and validation process,
//     the function may redirect the user, send error messages, or perform further authentication steps.
func OAuth2Callback(c *gin.Context) {
	span, ctx := utils.StartAPMSpan(c.Request.Context(), config.AppConfig.UseAPM, utils.GetFunctionName(), constants.OAUTH_ERRMSG_SPAN)
	defer utils.EndAPMSpan(span)

	code := c.DefaultQuery(constants.OAUTH_CODE_PARAM, "")
	if code == "" {
		c.HTML(http.StatusBadRequest, constants.UNAUTHORIZED_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_UNAUTHORIZED})
		return
	}

	token, err := config.AuthProvider.Exchange(c, code)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_EXCHANGE})
		return
	}

	if !token.Valid() {
		c.HTML(http.StatusInternalServerError, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_INVALID_TOKEN})
		return
	}

	sub, err := utils.GetSubjectFromToken(token)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constants.ERROR_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_SUB})
		return
	}

	utils.SetAPMUsername(span, ctx, sub)

	tgt := utils.GenerateTGT(config.AppConfig.TGTDuration, sub)
	setCookie(c, config.AppConfig.TGTName, tgt, config.AppConfig.Domain, config.AppConfig.TGTDuration)

	encryptedServiceURL, err := c.Cookie(constants.SERVICE_URL_COOKIE)
	unsetCookie(c, constants.SERVICE_URL_COOKIE, config.AppConfig.Domain)

	serviceURL, _ := utils.Decrypt(config.AppConfig.SecureCookie, encryptedServiceURL)
	if serviceURL != "" {
		action := getAction(serviceURL)
		if action != nil {
			action(c)
		}

		redirectToService(c, serviceURL, sub, tgt, true)
		return
	}

	c.HTML(http.StatusCreated, constants.LOGIN_HTML, gin.H{constants.TEMPLATE_MESSAGE: constants.OAUTH_ERRMSG_OK})
}

type Rule struct {
	ServiceURL string
	Action     string
	Params     map[string]string
}

var rules []Rule

func loadRules() []Rule {
	ruleStrings := strings.Split(config.AppConfig.Rules, "|")
	var loadedRules []Rule
	for _, ruleString := range ruleStrings {
		parts := strings.Split(ruleString, ";")
		if len(parts) >= 2 {
			serviceURL := parts[0]
			action := parts[1]
			params := make(map[string]string)
			for _, param := range parts[2:] {
				paramParts := strings.SplitN(param, "=", 2)
				if len(paramParts) == 2 {
					params[paramParts[0]] = paramParts[1]
				}
			}
			loadedRules = append(loadedRules, Rule{ServiceURL: serviceURL, Action: action, Params: params})
		}
	}
	return loadedRules
}

func getAction(serviceURL string) func(c *gin.Context) {
	rules = loadRules()

	parsedURL, err := url.Parse(serviceURL)
	if err != nil {
		log.Printf("Invalid serviceURL: %s, error: %v", serviceURL, err)
		return nil
	}
	domain := parsedURL.Hostname()
	log.Printf("Extracted domain: %s", domain)

	for _, rule := range rules {
		log.Printf("Rule: %+v", rule)
		if domain == rule.ServiceURL {
			log.Printf("Matched rule: %+v", rule)
			actionFunc := getActionFromEnv(rule.Action, rule.Params)
			if actionFunc != nil {
				return actionFunc
			}
		}
	}
	return nil
}

func getActionFromEnv(action string, params map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		switch action {
		case "unsetCookie":
			if cookieName, ok := params["cookieName"]; ok {
				unsetCookie(c, cookieName, config.AppConfig.Domain)
			}
		default:
		}
	}
}
