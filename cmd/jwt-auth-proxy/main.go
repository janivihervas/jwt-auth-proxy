package main

import (
	"github.com/gorilla/sessions"
	"github.com/janivihervas/authproxy"
	"github.com/janivihervas/authproxy/azure"
	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
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

	store := sessions.NewCookieStore([]byte(config.CookieHashKey), []byte(config.CookieEncryptKey))
	store.Options = &sessions.Options{
    Path: "/",
    MaxAge: int(time.Hour * 24 * 7),
    Secure:false,
    HttpOnly:true,
    SameSite: http.SameSiteDefaultMode,
	}

	m, err := authproxy.NewMiddleware(&authproxy.Config{
		AuthClient: &oauth2.Config{
			ClientID:     config.AzureClientID,
			ClientSecret: config.AzureClientSecret,
			Endpoint:     azure.Endpoint(config.AzureTenant),
			RedirectURL:  config.CallbackURL,
			Scopes:       []string{authproxy.ScopeOpenID, authproxy.ScopeOfflineAccess},
		},
		Next:                     upstream.Echo{},
		SessionStore:             store,
		SkipRedirectToLoginRegex: []string{"/api/.*"},
		Logger:                   log.New(os.Stdout, "", log.LstdFlags),
	})
	if err != nil {
		panic(err)
	}

	err = server.RunHTTP(config.Port, m)
	if err != nil {
		panic(err)
	}
}
