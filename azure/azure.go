package azure

import (
	"golang.org/x/oauth2"
)

// Endpoint will return Azure AD endpoint configuration for a given tenant id
func Endpoint(tenant string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/token",
	}
}
