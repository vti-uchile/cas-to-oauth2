package utils

import (
	"cas-to-oauth2/constants"
	"cas-to-oauth2/database"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
)

var (
	ticketExpiration          = 5 * time.Minute
	executionExpiration       = 1 * time.Minute
	executionCounter    int32 = 0
)

func GetSubjectFromToken(token *oauth2.Token) (string, error) {
	rawIDToken, ok := token.Extra(constants.UTILS_ID_TOKEN).(string)
	if !ok {
		return "", fmt.Errorf(constants.UTILS_ERRMSG_MISSING)
	}

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(rawIDToken, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf(constants.UTILS_ERRMSG_CLAIM_PARSE)
	}

	sub, ok := claims[constants.UTILS_CLAIM].(string)
	if !ok {
		return "", fmt.Errorf(constants.UTILS_ERRMSG_CLAIM_NOT_EXIST)
	}
	return sub, nil
}

func RandomString(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GenerateServiceTicket(service, username, tgt string, isDirect bool) string {
	st := fmt.Sprintf("ST-%s", RandomString(32))
	expiration := time.Now().Add(ticketExpiration)
	database.Conn.GenerateServiceTicket(st, service, username, isDirect, expiration)
	return st
}

func GenerateTGT(expire int, username string) string {
	tgt := fmt.Sprintf("TGT-%s", RandomString(32))
	timeMins := time.Duration(expire) * time.Minute
	expiration := time.Now().Add(timeMins)
	database.Conn.GenerateTGT(tgt, username, expiration)
	return tgt
}

func GenerateProxyTicket(service, pgt string) string {
	pt := fmt.Sprintf("PT-%s", RandomString(32))
	// expiration := time.Now().Add(ticketExpiration)
	// database.Conn.GenerateProxyTicket(pt, service, pgt, expiration)
	return pt
}

func ValidateServiceTicket(st, service string) (bool, string, bool) {
	return database.Conn.ValidateServiceTicket(st, service)
}

func ValidateTGT(tgt string) (bool, string) {
	return database.Conn.ValidateTGT(tgt)
}

func ValidatePGT(pgt string) (bool, error) {
	//return database.Conn.ValidatePGT(pgt)
	return true, nil
}

func ValidateProxyTicket(pt, service string) (string, string, []string, error) {
	//return database.Conn.ValidateProxyTicket(pt, service)
	return "", "", nil, nil
}

func DeleteTGT(tgt string) error {
	return database.Conn.DeleteTGT(tgt)
}

func IsTrue(s string) bool {
	return s == "true"
}

func Encrypt(secureCookie *securecookie.SecureCookie, value string) (string, error) {
	encoded, err := secureCookie.Encode(constants.SERVICE_URL_COOKIE, value)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

func Decrypt(secureCookie *securecookie.SecureCookie, value string) (string, error) {
	var decoded string
	err := secureCookie.Decode(constants.SERVICE_URL_COOKIE, value, &decoded)
	if err != nil {
		return "", err
	}
	return decoded, nil
}
