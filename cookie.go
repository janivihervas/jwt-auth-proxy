package authproxy

import (
	"net/http"
	"time"
)

func createAccessTokenCookie(accessToken string) *http.Cookie {
	return &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}
