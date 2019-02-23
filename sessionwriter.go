package authproxy

import (
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/sessions"
)

type sessionWriter struct {
	sessionStore sessions.Store
	session      *sessions.Session
	state        *sessionState
	logger       *log.Logger
	r            *http.Request
	http.ResponseWriter
}

func (s *sessionWriter) WriteHeader(statusCode int) {
	s.logger.Println(s.r.URL.Path, spew.Sdump(s.state))

	session, err := s.sessionStore.Get(s.r, sessionName)
	if err != nil {
		s.logger.Printf("sessionWriter: could not get session from request: %+v", err)
		s.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values[sessionName] = s.state
	err = session.Save(s.r, s)
	if err != nil {
		s.logger.Printf("sessionWriter: could not save session: %+v", err)
		s.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.ResponseWriter.WriteHeader(statusCode)
}
