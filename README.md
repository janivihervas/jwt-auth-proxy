# jwt-auth-proxy
Reverse proxy for JWT authentication

## Plan

Somewhat in order

- Directory structure
    - Go modules
- Simple reverse proxy
- CI + dev tools
    - [CircleCI](https://circleci.com/)
    - [CodeCov](https://codecov.io/)
    - [Code Climate](https://codeclimate.com/)?
    - [Go report](https://goreportcard.com/)
    - Issue tracking: Github milestones? [waffle.io](https://waffle.io/)?
- Write down plan in issues
- Kubernetes/Helm templates for deployment and testing
    - Actual Helm repository?
    - [Skaffold](https://github.com/GoogleContainerTools/skaffold)?
- Integration tests
    - [CrossBrowserTesting](https://crossbrowsertesting.com/)?
    - [BrowserStack](https://www.browserstack.com/)?
- Configuration
    - File?
    - Env vars?
    - Command line flags?
- Integration with [Azure AD](https://docs.microsoft.com/en-us/azure/active-directory/develop/v1-protocols-oauth-code):
    - [Get access token](https://docs.microsoft.com/en-us/azure/active-directory/develop/v1-protocols-oauth-code#oauth-20-authorization-flow)
    - Validate access token
    - How to pass access token? Header? Cookie?
- Repository: add PostgreSQL support
- Request refresh token and store in repository
    - Store cookie for refresh token
- Handle updating access token with refresh token
- Handle fetching a new a access token
    - Remember the url originally requested
    - Always redirect to Azure at this point
- Jaeger integration
- Configuration file
    - Which routes to protect?
    - API vs client mode (401 or redirect)
- Optional CSRF protection
- More repository providers, like SQLite (file), Redis, MySQL etc
