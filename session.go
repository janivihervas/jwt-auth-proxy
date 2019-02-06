package authproxy

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/pkg/errors"
)

func init() {
	gob.Register(State{})
}

const (
	sessionName = "authproxy_session"
)

// State of the current user
type State struct {
	AuthRequestState string
	AccessToken      string
	RefreshToken     string
	OriginalURL      string
}

func (m *Middleware) createNewSession(ctx context.Context, accessToken string, w http.ResponseWriter, r *http.Request) error {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return errors.Wrapf(err, "could not get session %s", sessionName)
	}

	// State is not set, but doesn't hurt to try to get it
	state, ok := session.Values[sessionName].(State)
	if !ok {
		// State is not stored, carry on
	}
	stateNew := State{
		AccessToken:      accessToken,
		RefreshToken:     state.RefreshToken,
		OriginalURL:      state.OriginalURL,
		AuthRequestState: state.AuthRequestState,
	}

	session.Values[sessionName] = stateNew
	err = session.Save(r, w)
	if err != nil {
		return errors.Wrapf(err, "couldn't save session %s", sessionName)
	}

	return nil
}

func (m *Middleware) clearSessionAndAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie := createAccessTokenCookie("")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	r.Header.Del(authHeaderName)

	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return errors.Wrapf(err, "couldn't get session %s", sessionName)
	}

	session.Values[sessionName] = State{}
	err = session.Save(r, w)
	if err != nil {
		return errors.Wrapf(err, "couldn't save session %s", sessionName)
	}

	return nil
}
