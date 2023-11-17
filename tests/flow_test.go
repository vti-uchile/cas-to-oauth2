package main

import (
	"cas-to-oauth2/config"
	"cas-to-oauth2/database"
	"cas-to-oauth2/internal/handlers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

var (
	originURL    string
	serviceURL   string
	testUser     string
	testPassword string
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../web/templates/*")

	config.LoadConfig()
	database.Connect(config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBURI,
		config.AppConfig.DBDatabase,
		config.AppConfig.DBPoolSize)

	fullUrl, err := url.Parse(os.Getenv("OAUTH2_REDIRECT_URL"))
	if err != nil {
		log.Fatalf("Error parsing the origin URL: %v", err)
	}

	originURL = fmt.Sprintf("%s://%s", fullUrl.Scheme, fullUrl.Host)
	serviceURL = os.Getenv("TEST_SERVICE_URL")
	testUser = os.Getenv("TEST_USER")
	testPassword = os.Getenv("TEST_PASSWORD")

	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.Login)
	r.GET("/oauth2/callback", handlers.OAuth2Callback)
	r.GET("/serviceValidate", handlers.ServiceValidate)
	r.GET("/validate", handlers.Validate)
	r.GET("/logout", handlers.Logout)
	r.POST("/logout", handlers.Logout)

	go func() {
		if err := r.Run(fullUrl.Host); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the test server: %v", err)
		}
	}()

	time.Sleep(time.Second)
	return r
}

func TestFlowHandler(t *testing.T) {
	setupTestRouter()

	var st string
	var err error
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	t.Run("GetTicketGrantingTicket", func(t *testing.T) {
		err := getCookie(client, originURL+"/login")
		if err != nil {
			t.Fatalf("Error obtaining the OAuth2 cookie: %s", err)
		}

		u, _ := url.Parse(originURL)
		cookies := client.Jar.Cookies(u)
		if len(cookies) == 0 {
			t.Error("Cookies were expected to be set, but none were found")
		}
		t.Logf("Cookies: %v", cookies)
	})

	t.Run("GetServiceTicket", func(t *testing.T) {
		st, err = getServiceTicket(client, originURL+"/login")
		if err != nil {
			t.Fatalf("Error getting the service ticket: %s", err)
		}
		if st == "" {
			t.Error("Expected a non-empty service ticket")
		}
		t.Logf("Service ticket: %s", st)
	})

	t.Run("ValidateServiceTicket", func(t *testing.T) {
		res, err := validateST(client, originURL+"/serviceValidate", st)
		if err != nil {
			t.Fatalf("Error validating the service ticket: %s", err)
		}
		t.Logf("Validate Response: %s", res)
	})

	t.Run("Logout", func(t *testing.T) {
		res, err := logout(client, originURL+"/logout")
		if err != nil {
			t.Fatalf("Error during logout: %s", err)
		}
		t.Logf("Logout Response: %s", res)
	})
}

func getCookie(client *http.Client, urlStr string) error {
	resp, err := client.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	actionURL, err := getFormData(resp)
	if err != nil {
		return err
	}

	loginData := url.Values{}
	loginData.Set("username", testUser)
	loginData.Set("password", testPassword)

	_, err = client.PostForm(actionURL, loginData)
	return err
}

func getFormData(resp *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	form := doc.Find("form").First()
	action, exists := form.Attr("action")
	if !exists {
		return "", fmt.Errorf("No action attribute found")
	}

	return action, nil
}

func getServiceTicket(client *http.Client, urlStr string) (string, error) {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("service", serviceURL)
	req.URL.RawQuery = q.Encode()

	_, resp, err := readAndCloseBody(client, req, false)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("a redirection was expected, status code received: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("'Location' header was not found for ticket extraction")
	}

	parsedURL, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	ticket := parsedURL.Query().Get("ticket")
	if ticket == "" {
		return "", fmt.Errorf("could not extract the 'ticket' from the URL")
	}

	return ticket, nil
}

func validateST(client *http.Client, validateURL, serviceTicket string) (string, error) {
	req, err := http.NewRequest("GET", validateURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("service", serviceURL)
	q.Add("ticket", serviceTicket)
	req.URL.RawQuery = q.Encode()

	body, resp, err := readAndCloseBody(client, req, true)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("a 200 status code was expected, but %d was received", resp.StatusCode)
	}

	return body, nil
}

func logout(client *http.Client, logoutURL string) (string, error) {
	req, err := http.NewRequest("GET", logoutURL, nil)
	if err != nil {
		return "", err
	}

	res, _, err := readAndCloseBody(client, req, true)
	if err != nil {
		return "", err
	}

	return res, nil
}

func readAndCloseBody(client *http.Client, req *http.Request, retBody bool) (string, *http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if !retBody {
		return "", resp, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	return string(body), resp, nil
}
