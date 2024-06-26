package config

import (
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/auth"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Config struct {
	OAuth2Server   string
	AuthMethod     string
	DBURI          string
	DBUser         string
	DBPassword     string
	DBDatabase     string
	DBPoolSize     int
	TGTName        string
	TGTDuration    int
	Domain         string
	AllowedDomains []string
	AnotherCookie  string
	JSessionID     string
	UseAPM         bool
	TGTSecure      bool
	TGTHttpOnly    bool
	SecureCookie   *securecookie.SecureCookie
}

var (
	AppConfig    Config
	AuthProvider auth.Authenticator
)

func LoadConfig() {
	if viper.GetString("DB_URI") == "" {
		godotenv.Load()
	}

	viper.AutomaticEnv()

	rawDomains := viper.GetString("ALLOWED_DOMAINS")

	AppConfig.OAuth2Server = viper.GetString("OAUTH2_SERVER")
	AppConfig.AuthMethod = viper.GetString("AUTH_METHOD")
	AppConfig.DBURI = viper.GetString("DB_URI")
	AppConfig.DBUser = viper.GetString("DB_USER")
	AppConfig.DBPassword = viper.GetString("DB_PASSWORD")
	AppConfig.DBDatabase = viper.GetString("DB_DATABASE")
	AppConfig.DBPoolSize, _ = strconv.Atoi(viper.GetString("DB_POOL_SIZE"))
	AppConfig.TGTName = viper.GetString("TGT_NAME")
	AppConfig.TGTDuration, _ = strconv.Atoi(viper.GetString("TGT_DURATION"))
	AppConfig.Domain = viper.GetString("DOMAIN_SCOPE")
	AppConfig.AllowedDomains = strings.Split(rawDomains, ",")
	AppConfig.AnotherCookie = viper.GetString("ANOTHER_COOKIE")
	AppConfig.JSessionID = viper.GetString("JSESSION_ID")
	AppConfig.UseAPM, _ = strconv.ParseBool(viper.GetString("USE_APM"))
	AppConfig.TGTSecure, _ = strconv.ParseBool(viper.GetString("TGT_SECURE"))
	AppConfig.TGTHttpOnly, _ = strconv.ParseBool(viper.GetString("TGT_HTTP_ONLY"))

	// SecureCookie is used to encrypt and decrypt the service URL
	hashKey := []byte(viper.GetString("SC_HASH_KEY"))
	blockKey := []byte(viper.GetString("SC_BLOCK_KEY"))
	AppConfig.SecureCookie = securecookie.New(hashKey, blockKey)

	if AppConfig.AuthMethod == constants.OAUTH_METHOD {
		AuthProvider = initOAuth2Provider()
	} else {
		log.Fatal("Auth method not supported")
	}
}

func initOAuth2Provider() auth.Authenticator {
	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("OAUTH2_CLIENT_ID"),
		ClientSecret: viper.GetString("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("OAUTH2_AUTH_URL"),
			TokenURL: viper.GetString("OAUTH2_TOKEN_URL"),
		},
	}

	return auth.NewOAuth2Authenticator(oauth2Config)
}
