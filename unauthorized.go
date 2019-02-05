package authproxy

//func (m *Middleware) unauthorized(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	err := m.clearSessionAndAccessToken(ctx, w, r)
//	if err != nil {
//		m.Logger.Printf("%+v", err)
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	ctx = r.Context()
//
//	state, err := m.getStateFromContext(ctx)
//	if err != nil {
//		m.Logger.Printf("%+v", err)
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	newState := random.String(32)
//	state.AuthRequestState = newState
//
//	err = m.SessionStore.Save(ctx, state.ID, state)
//	if err != nil {
//		m.Logger.Printf("%+v", err)
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	redirectURL := m.AuthClient.AuthCodeURL(state.AuthRequestState, oauth2.AccessTypeOffline)
//	for _, regexp := range m.skipRedirectToLoginRegex {
//		if regexp.MatchString(r.URL.Path) {
//				m.unauthorizedResponse(ctx, w, redirectURL)
//				return
//		}
//	}
//
//	state.OriginalURL = r.URL.String()
//	err = m.SessionStore.Save(ctx, state.ID, state)
//	if err != nil {
//		m.Logger.Printf("%+v", err)
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
//}
//
//type unauthorizedResponse struct {
//	StatusCode  int    `json:"statusCode"`
//	RedirectURL string `json:"redirectURL"`
//}
//
//func (m *Middleware) unauthorizedResponse(ctx context.Context, w http.ResponseWriter, redirectURL string) {
//	resp := unauthorizedResponse{
//		StatusCode:  http.StatusUnauthorized,
//		RedirectURL: redirectURL,
//	}
//	w.WriteHeader(http.StatusUnauthorized)
//
//	w.Header().Set("Content-Type", "application/json")
//	encoder := json.NewEncoder(w)
//	encoder.SetEscapeHTML(false)
//
//	err := encoder.Encode(resp)
//	if err != nil {
//		w.Header().Set("Content-Type", "text/plain")
//		_, err = w.Write([]byte(redirectURL))
//		if err != nil {
//			// log err
//		}
//	}
//}
