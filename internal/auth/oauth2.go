package auth

import (
	"cas-to-oauth2/internal/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OAuth2Authenticator struct {
	Config oauth2.Config
}

func NewOAuth2Authenticator(config oauth2.Config) *OAuth2Authenticator {
	return &OAuth2Authenticator{Config: config}
}

func (o *OAuth2Authenticator) Authenticate(username, password string) bool {
	token, err := o.Config.PasswordCredentialsToken(context.Background(), username, password)
	if err != nil {
		return false
	}

	return token.Valid()
}

func (o *OAuth2Authenticator) Exchange(c *gin.Context, code string) (*oauth2.Token, error) {
	return o.Config.Exchange(c, code)
}

func (o *OAuth2Authenticator) RedirectAuth(c *gin.Context) {
	state := utils.RandomString(32)
	authURL := o.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, authURL)
}
