package middleware

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/janivihervas/oidc-go/internal/random"

	"github.com/janivihervas/oidc-go"

	"github.com/janivihervas/oidc-go/jwt"
)

func (m *middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("default")
	b, _ := httputil.DumpRequest(r, true)
	log.Println(string(b))

	accessTokenStr, err := extractAccessToken(r)
	if err != nil {
		m.redirectToLogin(w, r)
		return
	}

	_, err = jwt.ParseAccessToken(accessTokenStr)
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
	log.Println("authorize callback")
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

	// TODO: Redirect to original url from session

	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
}

func (m *middleware) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("redirect to login")
	state := string(random.String(32))
	sessionID, session, err := m.session(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	session.State = state
	err = m.sessionStorage.Save(r.Context(), sessionID, session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, createSessionCookie(sessionID))

	url := m.client.AuthenticationRequestURL(oidc.ResponseModeFormPost, state)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (m *middleware) refreshAccessToken(writer http.ResponseWriter, r *http.Request) error {
	return nil
}
