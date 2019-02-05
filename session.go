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

//
//func (m *Middleware) getStateFromSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (session.State, error) {
//	var (
//		sessionID []byte
//	)
//
//	cookie, err := r.Cookie(sessionName)
//	if err != nil {
//		return session.State{}, errors.Wrap(err, "middleware: couldn't get session id cookie")
//	}
//
//	err = m.cookieStore.Decode(sessionName, cookie.Value, &sessionID)
//	if err != nil {
//		return session.State{}, errors.Wrap(err, "middleware: couldn't decode session id from cookie")
//	}
//
//	state, err := m.SessionStore.Get(ctx, sessionID)
//	if err != nil {
//		return session.State{}, errors.Wrap(err, "middleware: couldn't get session from storage")
//	}
//
//	return state, nil
//}

func (m *Middleware) createNewSession(ctx context.Context, accessToken string, w http.ResponseWriter, r *http.Request) error {
	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		return errors.Wrap(err, "could not get session")
	}

	state, ok := session.Values[sessionName].(State)
	if !ok {
		return errors.New("couldn't type cast")
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
		return errors.Wrap(err, "couldn't save session")
	}

	return nil
}

//
//func (m *Middleware) clearSessionAndAccessToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	cookie := createAccessTokenCookie("")
//	cookie.MaxAge = -1
//	http.SetCookie(w, cookie)
//
//	r.Header.Del(authHeaderName)
//
//	state, err := m.getStateFromContext(ctx)
//	if err == nil {
//		err = m.SessionStore.Delete(ctx, state.ID)
//		if err != nil {
//			m.Logger.Printf("%+v", err)
//		}
//	} else {
//		m.Logger.Printf("%+v", err)
//	}
//
//	err = m.createNewSession(ctx, "", w, r)
//	if err != nil {
//		return errors.Wrap(err, "middleware: couldn't create new session")
//	}
//
//	return nil
//}
