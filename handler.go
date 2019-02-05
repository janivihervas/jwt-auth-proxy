package authproxy

import (
	"net/http"
)

func (m *Middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	//for _, regexp := range m.skipAuthenticationRegex {
	//	if regexp.MatchString(r.URL.Path) {
	//		m.Next.ServeHTTP(w, r)
	//		return
	//	}
	//}
	//
	//var (
	//	accessTokenStr string
	//	ctx            = r.Context()
	//)
	//
	//accessTokenStr, err := m.getAccessToken(ctx, r)
	//if err != nil {
	//	accessTokenStr = m.refreshAccessToken(ctx, w, r)
	//}
	//
	//var validationErr error
	//for i := 0; i < 3; i++ {
	//	_, validationErr = jwt.ParseAccessToken(ctx, accessTokenStr)
	//	if validationErr == jwt.ErrTokenExpired {
	//		accessTokenStr = m.refreshAccessToken(ctx, w, r)
	//		continue
	//	}
	//	break
	//}
	//
	//if validationErr != nil {
	//	m.unauthorized(ctx, w, r)
	//	return
	//}
	//
	//m.Next.ServeHTTP(w, r)
}

func (m *Middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodGet {
	//	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//	return
	//}
	//
	//var (
	//	ctx      = r.Context()
	//	code     = r.URL.Query().Get("code")
	//	stateNew = r.URL.Query().Get("state")
	//)
	//
	//if code == "" {
	//	http.Error(w, "no code in response", http.StatusBadRequest)
	//	return
	//}
	//
	//if stateNew == "" {
	//	http.Error(w, "no state in response", http.StatusBadRequest)
	//	return
	//}
	//
	//state, err := m.getStateFromContext(ctx)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//stateOld := state.AuthRequestState
	//
	//// Clean previous state from session
	//state.AuthRequestState = ""
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//if stateNew != stateOld {
	//	http.Error(w, fmt.Sprintf("states are not the same: %s != %s", stateNew, stateOld), http.StatusBadRequest)
	//	return
	//}
	//
	//tokens, err := m.AuthClient.Exchange(ctx, code, oauth2.AccessTypeOffline)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//	return
	//}
	//
	//state.RefreshToken = tokens.RefreshToken
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))
	//
	//url := state.OriginalURL
	//
	//// Clear previous redirect url from session
	//state.OriginalURL = ""
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//}
	//
	//if url == "" {
	//	url = "/"
	//}
	//
	//http.Redirect(w, r, url, http.StatusSeeOther)
}
