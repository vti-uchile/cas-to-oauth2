package main

import (
	"bytes"
	"cas-to-oauth2/config"
	"cas-to-oauth2/database"
	"cas-to-oauth2/internal/handlers"
	"cas-to-oauth2/internal/utils"
	"fmt"
	"io"
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
	r.POST("/samlValidate", handlers.SamlValidate)
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

func TestStartServer(t *testing.T) {
	setupTestRouter()
}

func TestFlowHandler(t *testing.T) {
	var st string
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	t.Run("GetTicketGrantingTicket", func(t *testing.T) {
		getTicketGrantingTicket(t, client, originURL)
	})

	t.Run("GetServiceTicket", func(t *testing.T) {
		st = getServiceTicketTest(t, client, originURL+"/login")
	})

	t.Run("ValidateServiceTicket", func(t *testing.T) {
		res, err := validateST(client, originURL+"/serviceValidate", st)
		if err != nil {
			t.Fatalf("Error validating the service ticket: %s", err)
		}
		t.Logf("Validate Response: %s", res)
	})

	t.Run("Logout", func(t *testing.T) {
		logoutTest(t, client, originURL+"/logout")
	})
}

func TestFlowSAMLHandler(t *testing.T) {
	var st string
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	t.Run("GetTicketGrantingTicket", func(t *testing.T) {
		getTicketGrantingTicket(t, client, originURL)
	})

	t.Run("GetServiceTicket", func(t *testing.T) {
		st = getServiceTicketTest(t, client, originURL+"/login")
	})

	t.Run("ValidateSAMLTicket", func(t *testing.T) {
		res, err := validateSAML(client, originURL+"/samlValidate", st)
		if err != nil {
			t.Fatalf("Error validating the service ticket: %s", err)
		}
		t.Logf("Validate SAML Response: %s", res)
	})

	t.Run("Logout", func(t *testing.T) {
		logoutTest(t, client, originURL+"/logout")
	})
}

func getTicketGrantingTicket(t *testing.T, client *http.Client, originURL string) {
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
}

func getServiceTicketTest(t *testing.T, client *http.Client, loginURL string) string {
	st, err := getServiceTicket(client, loginURL)
	if err != nil {
		t.Fatalf("Error getting the service ticket: %s", err)
	}
	if st == "" {
		t.Error("Expected a non-empty service ticket")
	}
	t.Logf("Service ticket: %s", st)

	return st
}

func logoutTest(t *testing.T, client *http.Client, logoutURL string) {
	res, err := logout(client, logoutURL)
	if err != nil {
		t.Fatalf("Error during logout: %s", err)
	}
	t.Logf("Logout Response: %s", res)
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

func validateSAML(client *http.Client, validateURL, serviceTicket string) (string, error) {
	timeNow := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	randomID := utils.RandomString(43)

	xmlData := fmt.Sprintf(`<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
	<Header/>
	<Body>
		<samlp:Request xmlns:samlp="urn:oasis:names:tc:SAML:1.0:protocol" MajorVersion="1" MinorVersion="1" RequestID="_%s" IssueInstant="%s">
			<samlp:AssertionArtifact>%s</samlp:AssertionArtifact>
		</samlp:Request>
	</Body>
	</Envelope>`, randomID, timeNow, serviceTicket)

	req, err := http.NewRequest("POST", validateURL, bytes.NewBufferString(xmlData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/xml")
	q := req.URL.Query()
	q.Add("TARGET", serviceURL)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	return string(body), resp, nil
}
