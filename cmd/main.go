package main

import (
	"context"
	"log"
	"time"

	"cas-to-oauth2/config"
	"cas-to-oauth2/constants"
	"cas-to-oauth2/database"
	"cas-to-oauth2/internal/handlers"

	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
)

func main() {
	r := gin.Default()
	r.Use(apmgin.Middleware(r))

	r.LoadHTMLGlob("web/templates/*")
	config.LoadConfig()
	database.Connect(config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBURI,
		config.AppConfig.DBDatabase,
		config.AppConfig.DBPoolSize)

	checker := health.NewChecker(
		health.WithCheck(health.Check{
			Name: "mongodb",
			Check: func(ctx context.Context) error {
				return database.Conn.Client().Ping(ctx, nil)
			},
			Timeout: time.Second * 5,
		}),
	)

	r.GET(constants.ENDPOINT_ROOT, handlers.Login)
	r.GET(constants.ENDPOINT_LOGIN, handlers.Login)
	r.POST(constants.ENDPOINT_LOGIN, handlers.Login)
	r.GET(constants.ENDPOINT_OAUTH2, handlers.OAuth2Callback)
	r.GET(constants.ENDPOINT_SERVICE_VALIDATE, handlers.ServiceValidate)
	r.GET(constants.ENDPOINT_PROXY_VALIDATE, handlers.ServiceValidate)
	r.POST(constants.ENDPOINT_SAML_VALIDATE, handlers.SamlValidate)
	r.GET(constants.ENDPOINT_VALIDATE, handlers.Validate)
	r.GET(constants.ENDPOINT_LOGOUT, handlers.Proxy)
	r.GET(constants.ENDPOINT_LOGOUT, handlers.Logout)
	r.POST(constants.ENDPOINT_LOGOUT, handlers.Logout)
	r.GET(constants.ENDPOINT_HEALTHCHECK, gin.WrapF(health.NewHandler(checker)))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(constants.MAIN_ERRMSG, err)
	}
}
