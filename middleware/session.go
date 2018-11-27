package middleware

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/janivihervas/jwt-auth-proxy"
)

func (m *middleware) session(r *http.Request) ([]byte, oidc.Session, error) {
	var (
		sessionID []byte
		session   oidc.Session
	)

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return sessionID, session, errors.Wrapf(err, "middleware: couldn't get %s cookie", sessionCookieName)
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
