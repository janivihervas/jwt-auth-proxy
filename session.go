package authproxy

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/pkg/errors"
)

func init() {
	gob.Register(&sessionState{})
}

const (
	sessionName                = "authproxy_session"
	authRequestStateExpiration = time.Minute * 5
)

type authRequestState struct {
	ExpiresAt   time.Time
	OriginalURL string
}

// sessionState of the current request
type sessionState struct {
	// AuthRequestStates is stored for comparing the state returned by the authentication provider.
	AuthRequestStates map[string]authRequestState
	// AccessToken is stored so that authentication still works if the access token cookie is empty
	AccessToken string
	// RefreshToken for refreshing access token
	RefreshToken string
}

func (state *sessionState) clearExpiredStates() {
	now := time.Now()

	for k, v := range state.AuthRequestStates {
		if now.After(v.ExpiresAt) {
			delete(state.AuthRequestStates, k)
		}
	}
}

func (m *Middleware) initializeSession(r *http.Request) (*sessions.Session, *sessionState, error) {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get session %s", sessionName)
	}

	state, ok := session.Values[sessionName].(*sessionState)
	if !ok {
		state = &sessionState{
			AuthRequestStates: make(map[string]authRequestState),
		}
	}

	return session, state, nil
}

func (m *Middleware) clearSessionAndAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie := createAccessTokenCookie("")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	r.Header.Del(authHeaderName)

	state, err := getStateFromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get session from context")
	}
	state.AccessToken = ""
	// Do not clear state.AuthRequestStates
	state.RefreshToken = ""
	// Do not clear state.OriginalURLs

	return nil
}
