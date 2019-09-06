package main

type config struct {
	ConfigFile   string       `id:"config" short:"c" desc:"Configuration file"`
	AuthProxy    authProxy    `id:"auth-proxy" desc:"Settings for authproxy"`
	OIDC         oidc         `id:"oidc" desc:"Settings for OpenID Connect provider"`
	SessionStore sessionStore `id:"session-store" desc:"Settings for session store"`
}

type authProxy struct {
	Host                    string      `id:"host" desc:"Host to run the server in" default:"localhost"`
	Port                    int         `id:"port" desc:"Port to run the server in" default:"3000"`
	LivenessProbePort       int         `id:"liveness-probe-port" desc:"Port to run liveness probe. Will respond 200 when the server is running. Will respond the same to all paths" default:"3001"`
	ReadinessProbePort      int         `id:"readiness-probe-port" desc:"Port to run readiness probe. Will respond 200 when the server is healthy. Will respond the same to all paths" default:"3002"`
	CallbackPath            string      `id:"callback-path" desc:"Path for listening to authentication callback"`
	AuthenticationMode      string      `id:"authentication-mode" desc:"Authentication mode of the proxy. Accepted values are whitelist or blacklist. Use whitelist for authenticating only the urls specified in authentication-whitelist flag or blacklist for authenticating every url except the ones specified with authentication-blacklist flag. Will default to whitelist" default:"whitelist"`
	AuthenticationWhitelist []string    `id:"authentication-whitelist" desc:"Whitelisted regular expressions for paths to authenticate" default:".*"`
	AuthenticationBlacklist []string    `id:"authentication-blacklist" desc:"Blacklisted regular expressions for paths to not authenticate"`
	RedirectToLoginMode     string      `id:"redirect-to-login-mode" desc:"Login page redirection mode of the proxy. Accepted values are whitelist, blacklist or html-only. Determines whether to redirect an unauthenticated request or not. Non-interactive requests should not be redirected to login page. Use whitelist for redirecting only the urls specified in redirect-to-login-whitelist flag, blacklist for not redirecting urls specified in redirect-to-login-blacklist flag or html-only for redirecting only requests that have Accept: text/html header." default:"html-only"`
	AccessToken             accessToken `id:"access-token" desc:"Access token cookie settings"`
}

type accessToken struct {
	EnableCookie bool   `id:"enable-cookie" desc:"Whether to send the access token as a cookie" default:"true"`
	Name         string `id:"name" desc:"Cookie name" default:"access_token"`
	Path         string `id:"path" desc:"Cookie path" default:"/"`
	Domain       string `id:"domain" desc:"Cookie domain" default:""`
	HttpOnly     bool   `id:"http-only" desc:"True if cookie should not be accessible from Javascript" default:"true"`
	Secure       bool   `id:"secure" desc:"True if cookie is sent only over HTTPS" default:"true"`
	SameSite     string `id:"same-site" desc:"Cookie setting for SameSite. Accepted values are strict, lax or empty string" default:""`
}

type oidc struct {
	ClientID                              string     `id:"oidc.client-id" desc:"OIDC client id"`
	ClientSecret                          string     `id:"client-secret" desc:"OIDC client secret"`
	ConfigurationURL                      string     `id:"configuration-url" desc:"OIDC configuration url. Example: <issuer>/.well-known/openid-configuration"`
	CallbackURL                           string     `id:"callback-url" desc:"Full callback url to authproxy. Example: https://www.example.com/auth-callback"`
	AdditionalScopes                      []string   `id:"additional-scopes" desc:"Additional OIDC scopes. openid and offline_access are used by default"`
	AdditionalAuthenticationURLParameters []keyValue `id:"additional-authentication-url-parameters" desc:"Key-value pairs for providers who require additional authorization parameters. For instance Auth0 requires to 'audience' parameter to be set"`
}

type keyValue struct {
	Key   string `id:"key" desc:"Key"`
	Value string `id:"value" desc:"Value"`
}

type sessionStore struct {
	Backend string `id:"backend" desc:"Which backend to use for storing sessions. Only accepted value is redis"`
	Redis   redis  `id:"redis" desc:"Configuration of Redis"`
	Cookie  cookie `id:"cookie" desc:"Cookie settings"`
}

type redis struct {
	MaxIdleConnections int    `id:"max-idle-connections" desc:"Maximum idle connections" default:"100"`
	Network            string `id:"network" desc:"Network scheme" default:"tcp"`
	Address            string `id:"address" desc:"Host and port" default:"localhost:6379"`
	Password           string `id:"password" desc:"Password" default:""`
}

type cookie struct {
	Name     string `id:"name" desc:"Cookie name" default:"authproxy_session"`
	Path     string `id:"path" desc:"Cookie path" default:"/"`
	Domain   string `id:"domain" desc:"Cookie domain" default:""`
	HttpOnly bool   `id:"http-only" desc:"True if cookie should not be accessible from Javascript" default:"true"`
	Secure   bool   `id:"secure" desc:"True if cookie is sent only over HTTPS" default:"true"`
	SameSite string `id:"same-site" desc:"Cookie setting for SameSite. Accepted values are strict, lax or empty string" default:""`
	MaxAge   int    `id:"max-age" desc:"Max age in seconds. Defaults to one week" default:"604800"`
}
