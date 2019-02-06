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
	// AuthRequestState is stored for comparing the state returned by the authentication provider
	AuthRequestState string
	// AccessToken is stored so that authentication still works if the access token cookie is empty
	AccessToken string
	// RefreshToken for refreshing access token
	RefreshToken string
	// OriginalURL is the url that was requested but got redirected to login page.
	// When handling the authorization callback, middleware will redirect back to this url,
	// if found in session
	OriginalURL string
}

func (m *Middleware) createNewSession(ctx context.Context, accessToken string, w http.ResponseWriter, r *http.Request) error {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return errors.Wrapf(err, "could not get session %s", sessionName)
	}

	// State is not set, but doesn't hurt to try to get it
	state, _ := session.Values[sessionName].(State)

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
