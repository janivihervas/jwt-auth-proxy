package authproxy

import (
	"net/http"
)

func (m *Middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	m.Next.ServeHTTP(w, r)
}

func (m *Middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
}
