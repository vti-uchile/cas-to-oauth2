package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Authenticator interface {
	Authenticate(username, password string) bool
	RedirectAuth(c *gin.Context)
	Exchange(c *gin.Context, code string) (*oauth2.Token, error)
}
