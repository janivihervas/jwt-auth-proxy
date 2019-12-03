package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boj/redistore"

	"github.com/janivihervas/authproxy/azure"

	"golang.org/x/oauth2"

	"github.com/gorilla/handlers"

	"github.com/gorilla/sessions"
	"github.com/janivihervas/authproxy"
	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

type config struct {
	AzureClientID               string `envconfig:"AZURE_AD_CLIENT_ID"`
	AzureClientSecret           string `envconfig:"AZURE_AD_CLIENT_SECRET"`
	AzureTenant                 string `envconfig:"AZURE_AD_TENANT"`
	CookieHashKey               string `envconfig:"COOKIE_HASH_KEY"`
	CookieEncryptKey            string `envconfig:"COOKIE_ENCRYPT_KEY"`
	CallbackURL                 string `envconfig:"CALLBACK_URL"`
	CallbackPath                string `envconfig:"CALLBACK_PATH"`
	Port                        string `envconfig:"PORT"`
	AdditionalAuthURLParameters string `envconfig:"ADDITIONAL_AUTH_URL_PARAMETERS"`
}

func main() {
	var conf config

	_ = gotenv.OverLoad(".env")
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	store, err := redistore.NewRediStore(100, "tcp", ":6379", "", []byte(conf.CookieHashKey), []byte(conf.CookieEncryptKey))
	if err != nil {
		panic(err)
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(time.Hour * 24 * 7 / time.Second),
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
		AdditionalAuthURLParameters: parseAdditionalParameters(conf.AdditionalAuthURLParameters),
		Next:                        upstream.Echo{},
		CallbackPath:                conf.CallbackPath,
		SessionStore:                store,
		SkipAuthenticationRegex:     []string{"/static/.*"},
		SkipRedirectToLoginRegex:    []string{"/api/.*"},
		Logger:                      log.New(os.Stdout, "", log.LstdFlags),
	})
	if err != nil {
		panic(err)
	}

	err = server.RunHTTP(conf.Port, handlers.LoggingHandler(os.Stdout, m))
	if err != nil {
		panic(err)
	}
}

func parseAdditionalParameters(raw string) []oauth2.AuthCodeOption {
	var additionalParameters []oauth2.AuthCodeOption
	parameters := strings.Split(raw, "&")
	for _, parameter := range parameters {
		keyValue := strings.Split(parameter, "=")
		if len(keyValue) == 2 && keyValue[0] != "" && keyValue[1] != "" {
			additionalParameters = append(additionalParameters, oauth2.SetAuthURLParam(keyValue[0], keyValue[1]))
		}
	}

	return additionalParameters
}
