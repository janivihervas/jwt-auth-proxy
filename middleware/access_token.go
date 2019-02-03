package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/janivihervas/oidc-go/internal/random"
	"github.com/janivihervas/oidc-go/session"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"
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

	state, err := m.sessionStorage.Get(ctx, sessionID)
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

		err = m.sessionStorage.Save(ctx, state.ID, state)
		if err != nil {
			return errors.Wrap(err, "middleware: couldn't save session to storage")
		}

		sessionCookie := createSessionCookie(value)
		http.SetCookie(w, sessionCookie)
		r.AddCookie(sessionCookie)
	}

	return nil
}

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

func (m *Middleware) refreshAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	session, err := m.getSession(ctx, w, r, true)
	if err != nil {
		return ""
	}

	if session.RefreshToken == "" {
		return ""
	}

	accessToken, _ := m.getAccessToken(ctx, r)

	tokens, err := m.authClient.TokenSource(r.Context(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}).Token()
	if err != nil {
		return ""
	}

	if tokens.RefreshToken != "" {
		session.RefreshToken = tokens.RefreshToken
	}

	http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))
	err = m.sessionStorage.Save(ctx, session.ID, session)
	if err != nil {
		// log error
	}

	return tokens.AccessToken
}
