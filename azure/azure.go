package azure

import (
	"github.com/janivihervas/oidc-go"
)

// Endpoint will return Azure AD endpoint configuration for a given tenant id
func Endpoint(tenant string) oidc.Endpoint {
	return oidc.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/token",
	}
}
