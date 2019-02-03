package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
)

func extractAccessToken(ctx context.Context, r *http.Request) (string, error) {
	cookie, err := r.Cookie(accessTokenName)
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	header := r.Header.Get(authHeaderName)
	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == authHeaderPrefix && parts[1] != "" {
		return parts[1], nil
	}

	return "", errors.New("middleware: no access token in header or cookie")
}

func (m *middleware) refreshAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	session, err := m.getSession(ctx, w, r, true)
	if err != nil {
		return ""
	}

	if session.RefreshToken == "" {
		return ""
	}

	accessToken, _ := extractAccessToken(ctx, r)

	tokens, err := m.client.TokenSource(r.Context(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}).Token()
	spew.Dump(tokens)
	//tokens, err := m.client.RefreshTokens(ctx, session.RefreshToken)
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
