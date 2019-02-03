package main

import (
	"os"

	"golang.org/x/oauth2"

	"github.com/janivihervas/jwt-auth-proxy/session/memory"

	oidc "github.com/janivihervas/jwt-auth-proxy"
	"github.com/janivihervas/jwt-auth-proxy/azure"

	"github.com/janivihervas/jwt-auth-proxy/middleware"

	"github.com/janivihervas/jwt-auth-proxy/internal/server"

	"github.com/janivihervas/jwt-auth-proxy/upstream"

	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

type Config struct {
	AzureClientID     string `envconfig:"AZURE_AD_CLIENT_ID"`
	AzureClientSecret string `envconfig:"AZURE_AD_CLIENT_SECRET"`
	AzureTenant       string `envconfig:"AZURE_AD_TENANT"`
}

func main() {
	var config Config

	_ = gotenv.OverLoad(".env")
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	m := middleware.NewMiddleware(&oauth2.Config{
		ClientID:     config.AzureClientID,
		ClientSecret: config.AzureClientSecret,
		Endpoint:     azure.Endpoint(config.AzureTenant),
		RedirectURL:  "http://localhost:3000/oauth2/callback",
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	},
		memory.New(),
		upstream.Echo{},
	)

	err = server.RunHTTP(os.Getenv("PORT"), m)
	if err != nil {
		panic(err)
	}
}
