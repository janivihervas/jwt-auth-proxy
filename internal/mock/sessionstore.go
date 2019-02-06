package mock

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionStore mocks github.com/gorilla/sessions.Store
type SessionStore struct {
	ErrGet  error
	ErrNew  error
	ErrSave error
	Session *sessions.Session
}

// Get mock
func (store *SessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return store.Session, store.ErrGet
}

// New mock
func (store *SessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return store.Session, store.ErrNew
}

// Save mock
func (store *SessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return store.ErrSave
}
