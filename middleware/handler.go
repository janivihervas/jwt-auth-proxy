package middleware

import (
	"net/http"

	"github.com/janivihervas/oidc-go/jwt"
)

func (m *middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenStr, err := extractAccessToken(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessToken, err := jwt.ParseAccessToken(accessTokenStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = accessToken.Valid()
	if err == jwt.ErrTokenExpired {
		err = m.refreshAccessToken(w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			m.next.ServeHTTP(w, r)
		}
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	m.next.ServeHTTP(w, r)
}

func (m *middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response, err := m.client.AuthenticationResponseForm(r.Form)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sessionID, session, err := m.session(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer func() {
		err := m.sessionStorage.Save(r.Context(), sessionID, session)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	state := session.State
	session.State = ""

	if response.State != state {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tokens, err := m.client.TokenRequest(r.Context(), response.Code)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if tokens.RefreshToken != "" {
		session.RefreshToken = tokens.RefreshToken
	}

	http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))

	// TODO: What now? Redirect to somewhere
}
