package middleware

import (
	"log"
	"net/http"
	"time"
)

func createAccessTokenCookie(accessToken string) *http.Cookie {
	log.Println("create access token cookie")
	return &http.Cookie{
		Name:     accessTokenName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}

func createSessionCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}
