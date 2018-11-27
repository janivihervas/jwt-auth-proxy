package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

func extractAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(accessTokenName)
	if err != nil || cookie.Value == "" {
		return cookie.Value, nil
	}

	header := r.Header.Get(authHeaderName)
	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == authHeaderPrefix && parts[1] != "" {
		return parts[1], nil
	}

	return "", errors.New("middleware: no access token in header or cookie")
}

func createAccessTokenCookie(accessToken string) *http.Cookie {
	return &http.Cookie{
		Name:     accessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}
