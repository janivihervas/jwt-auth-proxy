package authproxy

import (
	"context"
	"net/http"
)

// Middleware for authentication requests
type Middleware struct {
	*Config
}

// ServeHTTP will authenticate the request and forward it to the next http.Handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, state, err := m.initializeSession(r)
	if err != nil {
		m.Logger.Printf("couldn't initialise session: %+v", err)
		http.Error(w, "couldn't initialise session", http.StatusInternalServerError)
		return
	}

	ctx = context.WithValue(ctx, ctxStateKey, state)
	r = r.WithContext(ctx)

	err = m.setupAccessToken(ctx, w, r)
	if err != nil {
		// Access token is not set in request
		m.Logger.Printf("couldn't setup access token: %+v", err)
	}

	writer := &sessionWriter{
		sessionStore:   m.SessionStore,
		session:        session,
		state:          state,
		logger:         m.Logger,
		r:              r,
		ResponseWriter: w,
	}
	m.mux.ServeHTTP(writer, r)
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(config *Config) (*Middleware, error) {
	if err := config.Valid(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	config.mux = mux

	m := &Middleware{
		Config: config,
	}
	mux.HandleFunc(m.CallbackPath, m.authorizeCallback)
	mux.HandleFunc("/", m.defaultHandler)

	return m, nil
}
