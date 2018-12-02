package middleware

import (
	"context"
	"net/http"

	"github.com/janivihervas/oidc-go"

	"github.com/janivihervas/oidc-go/internal/random"
	"github.com/pkg/errors"
)

func (m *middleware) getSession(ctx context.Context, w http.ResponseWriter, r *http.Request, createNew bool) (oidc.Session, error) {
	var (
		sessionID []byte
	)

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		// Cookie not found
		if createNew {
			return m.createNewSession(ctx, w)
		}
		return oidc.Session{}, errors.Wrap(err, "middleware: couldn't get session id cookie")
	}

	err = m.cookieStore.Decode(sessionCookieName, cookie.Value, &sessionID)
	if err != nil {
		// Something wrong with the cookie
		if createNew {
			return m.createNewSession(ctx, w)
		}
		return oidc.Session{}, errors.Wrap(err, "middleware: couldn't decode session id from cookie")
	}

	session, err := m.sessionStorage.Get(ctx, sessionID)
	if err != nil {
		// Session is not stored in storage
		if createNew {
			return m.createNewSession(ctx, w)
		}
		return oidc.Session{}, errors.Wrap(err, "middleware: couldn't get session from storage")
	}

	return session, nil
}

func (m *middleware) createNewSession(ctx context.Context, w http.ResponseWriter) (oidc.Session, error) {
	newSession := oidc.Session{
		ID: random.Bytes(32),
	}

	value, err := m.cookieStore.Encode(sessionCookieName, newSession.ID)
	if err != nil {
		return newSession, errors.Wrap(err, "middleware: couldn't encode session ID")
	}

	err = m.sessionStorage.Save(ctx, newSession.ID, newSession)
	if err != nil {
		return newSession, errors.Wrap(err, "middleware: couldn't save session to storage")
	}

	http.SetCookie(w, createSessionCookie(value))

	return newSession, nil
}

func (m *middleware) clearSessionAndAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cookie := createAccessTokenCookie("")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	cookie = createSessionCookie("")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	session, err := m.getSession(ctx, w, r, false)
	if err != nil {
		// log err
		return
	}

	err = m.sessionStorage.Delete(ctx, session.ID)
	if err != nil {
		// log err
	}
}
