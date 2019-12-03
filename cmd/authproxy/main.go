package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

func main() {
	var conf config

	_ = gotenv.OverLoad(".env")
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	//store, err := redistore.NewRediStore(100, "tcp", ":6379", "", []byte(conf.CookieHashKey), []byte(conf.CookieEncryptKey))
	//if err != nil {
	//	panic(err)
	//}
	//store.Options = &sessions.Options{
	//	Path:     "/",
	//	MaxAge:   int(time.Hour * 24 * 7 / time.Second),
	//	Secure:   false,
	//	HTTPOnly: true,
	//	SameSite: http.SameSiteDefaultMode,
	//}
	//
	//m, err := authproxy.NewMiddleware(&authproxy.Config{
	//	AuthClient: &oauth2.Config{
	//		ClientID:     conf.AzureClientID,
	//		ClientSecret: conf.AzureClientSecret,
	//		Endpoint:     azure.Endpoint(conf.AzureTenant),
	//		RedirectURL:  conf.CallbackURL,
	//		Scopes:       []string{authproxy.ScopeOpenID, authproxy.ScopeOfflineAccess},
	//	},
	//	AdditionalAuthURLParameters: parseAdditionalParameters(conf.AdditionalAuthURLParameters),
	//	Next:                        upstream.Echo{},
	//	CallbackPath:                conf.CallbackPath,
	//	SessionStore:                store,
	//	SkipAuthenticationRegex:     []string{"/static/.*"},
	//	SkipRedirectToLoginRegex:    []string{"/api/.*"},
	//	Logger:                      log.New(os.Stdout, "", log.LstdFlags),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = server.RunHTTP(conf.Port, handlers.LoggingHandler(os.Stdout, m))
	//if err != nil {
	//	panic(err)
	//}
}

//func parseAdditionalParameters(raw string) []oauth2.AuthCodeOption {
//	var additionalParameters []oauth2.AuthCodeOption
//	parameters := strings.Split(raw, "&")
//	for _, parameter := range parameters {
//		keyValue := strings.Split(parameter, "=")
//		if len(keyValue) == 2 && keyValue[0] != "" && keyValue[1] != "" {
//			additionalParameters = append(additionalParameters, oauth2.SetAuthURLParam(keyValue[0], keyValue[1]))
//		}
//	}
//
//	return additionalParameters
//}
