package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/janivihervas/oidc-go"
	"github.com/janivihervas/oidc-go/internal/random"
)

func extractAccessToken(r *http.Request) (string, error) {
	log.Println("extract access token")
	cookie, err := r.Cookie(accessTokenName)
	if err == nil && cookie.Value == "" {
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
	log.Println("create access token cookie")
	return &http.Cookie{
		Name:     accessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}

func createSessionCookie(sessionID []byte) *http.Cookie {
	log.Println("create session cookie")
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    string(sessionID),
		HttpOnly: true,
		Secure:   false,       // TODO
		Expires:  time.Time{}, // TODO
		MaxAge:   0,           // TODO
	}
}

func (m *middleware) session(r *http.Request) ([]byte, oidc.Session, error) {
	var (
		sessionID []byte
		session   oidc.Session
	)

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		session.ID = random.Byte(32)
		return session.ID, session, m.sessionStorage.Save(r.Context(), session.ID, session)
	}

	err = m.cookieStore.Decode(sessionCookieName, cookie.Value, &sessionID)
	if err != nil {
		return sessionID, session, errors.Wrapf(err, "middleware: couldn't decode %s cookie", sessionCookieName)
	}

	if len(sessionID) == 0 {
		return sessionID, session, errors.New("middleware: session id was empty")
	}

	session, err = m.sessionStorage.Get(r.Context(), sessionID)
	if err != nil {
		return sessionID, session, errors.Wrap(err, "middleware: couldn't session from storage")
	}

	return sessionID, session, nil
}
