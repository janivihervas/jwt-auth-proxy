package main

import (
	"os"

	"github.com/janivihervas/oidc-go/session/memory"

	"github.com/janivihervas/oidc-go"
	"github.com/janivihervas/oidc-go/azure"

	"github.com/janivihervas/oidc-go/middleware"

	"github.com/janivihervas/oidc-go/internal/server"

	"github.com/janivihervas/oidc-go/internal/http/upstream"

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

	m := middleware.New(oidc.Client{
		ClientID:     config.AzureClientID,
		ClientSecret: config.AzureClientSecret,
		Endpoint:     azure.Endpoint(config.AzureTenant),
		RedirectURL:  "http://localhost:3000/oauth2/callback",
		Scope:        []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	},
		memory.New(),
		upstream.Echo{},
	)

	err = server.RunHTTP(os.Getenv("PORT"), m)
	if err != nil {
		panic(err)
	}
}
