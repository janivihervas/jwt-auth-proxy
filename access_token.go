package authproxy

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/janivihervas/authproxy/internal/random"
	"github.com/janivihervas/authproxy/session"
)

func (m *Middleware) getAccessTokenFromCookie(ctx context.Context, r *http.Request) string {
	cookie, err := r.Cookie(accessTokenName)
	if err != nil || cookie.Value == "" {
		return ""
	}

	return cookie.Value
}

func (m *Middleware) getAccessTokenFromHeader(ctx context.Context, r *http.Request) string {
	header := r.Header.Get(authHeaderName)
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != authHeaderPrefix || parts[1] == "" {
		return ""
	}

	return parts[1]
}

func (m *Middleware) getAccessTokenFromSession(ctx context.Context, r *http.Request) string {
	var sessionID []byte

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}

	err = m.cookieStore.Decode(sessionCookieName, cookie.Value, &sessionID)
	if err != nil {
		return ""
	}

	state, err := m.SessionStorage.Get(ctx, sessionID)
	if err != nil || state.AccessToken == "" {
		return ""
	}

	return state.AccessToken
}

func (m *Middleware) setupAccessTokenAndSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var (
		cookieSet   bool
		headerSet   bool
		sessionSet  bool
		accessToken string
	)

	if s := m.getAccessTokenFromCookie(ctx, r); s != "" {
		accessToken = s
		cookieSet = true
	}

	if s := m.getAccessTokenFromHeader(ctx, r); s != "" {
		accessToken = s
		headerSet = true
	}

	if s := m.getAccessTokenFromSession(ctx, r); s != "" {
		accessToken = s
		sessionSet = true
	}

	if accessToken == "" {
		// TODO: what now?
		return nil
	}

	if !cookieSet {
		accessTokenCookie := createAccessTokenCookie(accessToken)
		http.SetCookie(w, accessTokenCookie)
		r.AddCookie(accessTokenCookie)
	}
	if !headerSet {
		r.Header.Set(authHeaderName, authHeaderPrefix+" "+accessToken)
	}
	if !sessionSet {
		state := session.State{
			ID:          random.Bytes(32),
			AccessToken: accessToken,
		}

		value, err := m.cookieStore.Encode(sessionCookieName, state.ID)
		if err != nil {
			return errors.Wrap(err, "middleware: couldn't encode session ID")
		}

		err = m.SessionStorage.Save(ctx, state.ID, state)
		if err != nil {
			return errors.Wrap(err, "middleware: couldn't save session to storage")
		}

		sessionCookie := createSessionCookie(value)
		http.SetCookie(w, sessionCookie)
		r.AddCookie(sessionCookie)
	}

	return nil
}

//lint:ignore U1000 not in use yet
func (m *Middleware) getAccessToken(ctx context.Context, r *http.Request) (string, error) {
	if s := m.getAccessTokenFromCookie(ctx, r); s != "" {
		return s, nil
	}

	if s := m.getAccessTokenFromHeader(ctx, r); s != "" {
		return s, nil
	}

	if s := m.getAccessTokenFromSession(ctx, r); s != "" {
		return s, nil
	}

	return "", errors.New("middleware: no access token in cookie, header or session")
}
