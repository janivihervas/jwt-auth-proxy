package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/boj/redistore"

	"github.com/gorilla/sessions"
	"github.com/janivihervas/authproxy"
	"github.com/janivihervas/authproxy/azure"
	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
)

type config struct {
	AzureClientID     string `envconfig:"AZURE_AD_CLIENT_ID"`
	AzureClientSecret string `envconfig:"AZURE_AD_CLIENT_SECRET"`
	AzureTenant       string `envconfig:"AZURE_AD_TENANT"`
	CookieHashKey     string `envconfig:"COOKIE_HASH_KEY"`
	CookieEncryptKey  string `envconfig:"COOKIE_ENCRYPT_KEY"`
	CallbackURL       string `envconfig:"CALLBACK_URL"`
	Port              string `envconfig:"PORT"`
}

func main() {
	var conf config

	_ = gotenv.OverLoad(".env")
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	//store := sessions.NewCookieStore([]byte(conf.CookieHashKey), []byte(conf.CookieEncryptKey))
	store, err := redistore.NewRediStore(100, "tcp", ":6379", "", []byte(conf.CookieHashKey), []byte(conf.CookieEncryptKey))
	if err != nil {
		panic(err)
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(time.Hour * 24 * 7),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	m, err := authproxy.NewMiddleware(&authproxy.Config{
		AuthClient: &oauth2.Config{
			ClientID:     conf.AzureClientID,
			ClientSecret: conf.AzureClientSecret,
			Endpoint:     azure.Endpoint(conf.AzureTenant),
			RedirectURL:  conf.CallbackURL,
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

	err = server.RunHTTP(conf.Port, handlers.LoggingHandler(os.Stdout, m))
	if err != nil {
		panic(err)
	}
}
