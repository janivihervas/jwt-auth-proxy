package provider

// OIDC provider
type OIDC interface {
	AuthorizationUrl() string
}
