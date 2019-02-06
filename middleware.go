package authproxy

import (
	"net/http"
)

// Middleware for authentication requests
type Middleware struct {
	*Config
}

// ServeHTTP will authenticate the request and forward it to the next http.Handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := m.setupAccessTokenAndSession(ctx, w, r)
	if err != nil {
		// Either access token is empty or couldn't create a session
		m.Logger.Printf("couldn't setup access token or session: %+v", err)
	}

	m.mux.ServeHTTP(w, r)
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
