package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func (m *middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("callback")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	response, err := m.client.ParseAuthenticationResponseForm(ctx, r.Form)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	session, err := m.getSession(ctx, w, r, true)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	state := session.State

	// Clean previous state from session
	session.State = ""
	err = m.sessionStorage.Save(ctx, session.ID, session)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
	}

	if response.State != state {
		//http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("%s != %s", response.State, state), http.StatusBadRequest)
		return
	}

	tokens, err := m.client.TokenRequest(ctx, response.Code)
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
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
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
