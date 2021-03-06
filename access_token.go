package authproxy

import (
	"context"
	"net/http"
	"strings"

	"github.com/janivihervas/authproxy/jwt"

	"golang.org/x/oauth2"

	"github.com/pkg/errors"
)

const (
	accessTokenCookieName = "access_token"
	authHeaderName        = "Authorization"
	authHeaderPrefix      = "Bearer"
)

func (m *Middleware) getAccessTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get cookie %s", accessTokenCookieName)
	}
	if cookie.Value == "" {
		return "", errors.Errorf("cookie %s is empty", accessTokenCookieName)
	}

	return cookie.Value, nil
}

func (m *Middleware) getAccessTokenFromHeader(r *http.Request) (string, error) {
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

func (m *Middleware) getAccessTokenFromSession(ctx context.Context) (string, error) {
	state, err := getStateFromContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get session from context")
	}

	if state.AccessToken == "" {
		return "", errors.New("access token in session is empty")
	}

	return state.AccessToken, nil
}

func (m *Middleware) setupAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var (
		cookieSet   bool
		headerSet   bool
		sessionSet  bool
		accessToken string
	)

	if s, err := m.getAccessTokenFromCookie(r); err == nil {
		accessToken = s
		cookieSet = true
	}

	if s, err := m.getAccessTokenFromHeader(r); err == nil {
		accessToken = s
		headerSet = true
	}

	if s, err := m.getAccessTokenFromSession(ctx); err == nil {
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

	if !sessionSet {
		state, err := getStateFromContext(ctx)
		if err != nil {
			return err
		}
		state.AccessToken = accessToken
	}

	return nil
}

func (m *Middleware) getAccessToken(ctx context.Context, r *http.Request) (string, error) {
	if s, err := m.getAccessTokenFromCookie(r); err == nil {
		return s, nil
	}

	if s, err := m.getAccessTokenFromHeader(r); err == nil {
		return s, nil
	}

	if s, err := m.getAccessTokenFromSession(ctx); err == nil {
		return s, nil
	}

	return "", errors.New("no access token in cookie, header or session")
}

func (m *Middleware) refreshAccessToken(ctx context.Context, w http.ResponseWriter) (string, error) {
	state, err := getStateFromContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get session from context")
	}

	if state.RefreshToken == "" {
		return "", errors.New("no refresh token in session")
	}

	tokens, err := m.AuthClient.TokenSource(ctx, &oauth2.Token{
		AccessToken:  state.AccessToken,
		RefreshToken: state.RefreshToken,
	}).Token()
	if err != nil {
		return "", errors.Wrap(err, "couldn't refresh tokens")
	}

	_, err = jwt.ParseAccessToken(ctx, tokens.AccessToken)
	if err != nil {
		return "", errors.Wrap(err, "couldn't refresh tokens, returned access token was invalid")
	}

	state.AccessToken = tokens.AccessToken
	http.SetCookie(w, createAccessTokenCookie(state.AccessToken))

	if tokens.RefreshToken != "" {
		state.RefreshToken = tokens.RefreshToken
	}

	return state.AccessToken, nil
}
