package middleware

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/davecgh/go-spew/spew"
)

func (m *middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("callback")
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var (
		code  = r.URL.Query().Get("code")
		state = r.URL.Query().Get("state")
	)

	if code == "" {
		http.Error(w, "no code in response", http.StatusBadRequest)
		return
	}

	if state == "" {
		http.Error(w, "no state in response", http.StatusBadRequest)
		return
	}

	session, err := m.getSession(ctx, w, r, true)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	stateOld := session.State

	// Clean previous state from session
	session.State = ""
	err = m.sessionStorage.Save(ctx, session.ID, session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if state != stateOld {
		http.Error(w, fmt.Sprintf("states are not the same: %s != %s", state, stateOld), http.StatusBadRequest)
		return
	}

	tokens, err := m.client.Exchange(ctx, code, oauth2.AccessTypeOffline)
	spew.Dump(tokens)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	if tokens.RefreshToken != "" {
		// Refresh refresh_token
		session.RefreshToken = tokens.RefreshToken
		err = m.sessionStorage.Save(ctx, session.ID, session)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))

	url := session.OriginalURL

	// Clear previous redirect url from session
	session.OriginalURL = ""
	err = m.sessionStorage.Save(ctx, session.ID, session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if url == "" {
		url = "/"
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
