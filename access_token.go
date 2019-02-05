package authproxy

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2"

	"github.com/pkg/errors"
)

const (
	accessTokenCookieName = "access_token"
	authHeaderName        = "Authorization"
	authHeaderPrefix      = "Bearer"
)

func (m *Middleware) getAccessTokenFromCookie(ctx context.Context, r *http.Request) (string, error) {
	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get cookie %s", accessTokenCookieName)
	}
	if cookie.Value == "" {
		return "", errors.Errorf("cookie %s is empty", accessTokenCookieName)
	}

	return cookie.Value, nil
}

func (m *Middleware) getAccessTokenFromHeader(ctx context.Context, r *http.Request) (string, error) {
	header := r.Header.Get(authHeaderName)
	if header == "" {
		return "", errors.Errorf("header %s is empty", authHeaderName)
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != authHeaderPrefix || parts[1] == "" {
		return "", errors.Errorf("header %s is malformed: %s", authHeaderName, header)
	}

	return parts[1], nil
}

func (m *Middleware) getAccessTokenFromSession(ctx context.Context, r *http.Request) (string, error) {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get session %s", sessionName)
	}

	state, ok := session.Values[sessionName].(State)
	if !ok {
		return "", errors.Errorf("couldn't type case session %s", sessionName)
	}

	return state.AccessToken, nil
}

func (m *Middleware) setupAccessTokenAndSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var (
		cookieSet   bool
		headerSet   bool
		sessionSet  bool
		accessToken string
	)

	if s, err := m.getAccessTokenFromCookie(ctx, r); err != nil {
		accessToken = s
		cookieSet = true
	}

	if s, err := m.getAccessTokenFromHeader(ctx, r); err != nil {
		accessToken = s
		headerSet = true
	}

	if s, err := m.getAccessTokenFromSession(ctx, r); err != nil {
		accessToken = s
		sessionSet = true
	}

	if accessToken == "" {
		return errors.New("access token is not set")
	}

	if !cookieSet {
		accessTokenCookie := createAccessTokenCookie(accessToken)
		http.SetCookie(w, accessTokenCookie)
		r.AddCookie(accessTokenCookie)
	}
	if !headerSet {
		r.Header.Set(authHeaderName, authHeaderPrefix+" "+accessToken)
	}

	// Always create the session so the next handlers don't need to do it
	if !sessionSet {
		return errors.Wrap(m.createNewSession(ctx, accessToken, w, r), "couldn't create new session")
	}

	return nil
}

func (m *Middleware) getAccessToken(ctx context.Context, r *http.Request) (string, error) {
	if s, err := m.getAccessTokenFromCookie(ctx, r); err != nil {
		return s, nil
	}

	if s, err := m.getAccessTokenFromHeader(ctx, r); err != nil {
		return s, nil
	}

	if s, err := m.getAccessTokenFromSession(ctx, r); err != nil {
		return s, nil
	}

	return "", errors.New("no access token in cookie, header or session")
}

func (m *Middleware) refreshAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get session %s", sessionName)
	}

	state, ok := session.Values[sessionName].(State)
	if !ok {
		return "", errors.Errorf("couldn't type cast session %s", err)
	}

	if state.RefreshToken == "" {
		return "", errors.New("no refresh token in session")
	}

	tokens, err := m.AuthClient.TokenSource(r.Context(), &oauth2.Token{
		AccessToken:  state.AccessToken,
		RefreshToken: state.RefreshToken,
	}).Token()
	if err != nil {
		return "", errors.Wrap(err, "couldn't refresh tokens")
	}

	if tokens.RefreshToken != "" {
		state.RefreshToken = tokens.RefreshToken
	}
	state.AccessToken = tokens.AccessToken

	session.Values[sessionName] = state
	err = session.Save(r, w)
	if err != nil {
		return state.AccessToken, errors.New("couldn't save session")
	}

	return state.AccessToken, nil
}
