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
	"github.com/subosito/gotenv"

	"github.com/stevenroose/gonfig"
)

type config struct {
	ConfigFile                  string `id:"config" short:"c" desc:"Configuration file"`
	ClientID                    string `id:"client-id" desc:"OIDC client id"`
	ClientSecret                string `id:"client-secret" desc:"OIDC client secret"`
	CookieHashKey               string `id:"cookie-hash-key" desc:"Key for hashing session cookie"`
	CookieEncryptKey            string `id:"cookie-encrypt-key" desc:"Key for encrypting session cookie"`
	CallbackURL                 string `id:"callback-url" desc:"Full callback url to authproxy. Example: https://www.example.com/auth-callback"`
	callbackPath                string
	Port                        string `id:"port" short:"p" default:"3000" desc:"Port to run the server on"`
	AdditionalAuthURLParameters string `id:"additional-auth-url-parameters" desc:"Key-value pairs for providers who require additional authorization parameters. Example: audience=auth0appName,foo=bar"`
}

func main() {
	var conf config
	_ = gotenv.OverLoad(".env")

	err := gonfig.Load(&conf, gonfig.Conf{})
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
