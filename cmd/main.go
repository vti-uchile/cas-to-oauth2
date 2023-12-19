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

	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.Login)
	r.GET("/oauth2/callback", handlers.OAuth2Callback)
	r.GET("/serviceValidate", handlers.ServiceValidate)
	r.POST("/samlValidate", handlers.SamlValidate)
	r.GET("/validate", handlers.Validate)
	r.GET("/proxy", handlers.Proxy)
	r.GET("/logout", handlers.Logout)
	r.POST("/logout", handlers.Logout)
	r.GET("/healthcheck", gin.WrapF(health.NewHandler(checker)))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(constants.MAIN_ERRMSG, err)
	}
}
