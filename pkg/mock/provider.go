package mock

// Provider mock
type Provider struct {
	ValidationError error
}

func (Provider) AuthorizationUrl() string {
	return "https://oidc.com/oauth2/authorize"
}
