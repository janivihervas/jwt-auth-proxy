package main

import (
	"github.com/janivihervas/authproxy"
	"github.com/janivihervas/authproxy/azure"
	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/session/memory"
	"github.com/janivihervas/authproxy/upstream"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
)

type Config struct {
	AzureClientID     string `envconfig:"AZURE_AD_CLIENT_ID"`
	AzureClientSecret string `envconfig:"AZURE_AD_CLIENT_SECRET"`
	AzureTenant       string `envconfig:"AZURE_AD_TENANT"`
	CookieHashKey string `envconfig:"COOKIE_HASH_KEY"`
	CookieEncryptKey string `envconfig:"COOKIE_ENCRYPT_KEY"`
	CallbackURL string `envconfig:"CALLBACK_URL"`
	Port string `envconfig:"PORT"`
}

func main() {
	var config Config

	_ = gotenv.OverLoad(".env")
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	m, err := authproxy.NewMiddleware(&authproxy.Config{
		AuthClient: &oauth2.Config{
			ClientID:     config.AzureClientID,
			ClientSecret: config.AzureClientSecret,
			Endpoint:     azure.Endpoint(config.AzureTenant),
			RedirectURL:  config.CallbackURL,
			Scopes:       []string{authproxy.ScopeOpenID, authproxy.ScopeOfflineAccess},
		},
		SessionStorage:             memory.New(),
		Next:                       upstream.Echo{},
		CookieHashKey:              []byte(config.CookieHashKey),
		CookieEncryptKey:           []byte(config.CookieEncryptKey),
		SkipRedirectToLoginRegex:   []string{"/api/.*"},
	})
	if err != nil {
		panic(err)
	}

	err = server.RunHTTP(config.Port, m)
	if err != nil {
		panic(err)
	}
}
