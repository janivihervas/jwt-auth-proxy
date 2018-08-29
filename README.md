# jwt-auth-proxy

[![CircleCI](https://circleci.com/gh/janivihervas/jwt-auth-proxy.svg?style=svg)](https://circleci.com/gh/janivihervas/jwt-auth-proxy)

[![Go Report Card](https://goreportcard.com/badge/github.com/janivihervas/jwt-auth-proxy)](https://goreportcard.com/report/github.com/janivihervas/jwt-auth-proxy)
[![GoDoc](https://godoc.org/github.com/janivihervas/jwt-auth-proxy?status.svg)](https://godoc.org/github.com/janivihervas/jwt-auth-proxy)
[![Sponsored](https://img.shields.io/badge/chilicorn-sponsored-brightgreen.svg?logo=data%3Aimage%2Fpng%3Bbase64%2CiVBORw0KGgoAAAANSUhEUgAAAA4AAAAPCAMAAADjyg5GAAABqlBMVEUAAAAzmTM3pEn%2FSTGhVSY4ZD43STdOXk5lSGAyhz41iz8xkz2HUCWFFhTFFRUzZDvbIB00Zzoyfj9zlHY0ZzmMfY0ydT0zjj92l3qjeR3dNSkoZp4ykEAzjT8ylUBlgj0yiT0ymECkwKjWqAyjuqcghpUykD%2BUQCKoQyAHb%2BgylkAyl0EynkEzmkA0mUA3mj86oUg7oUo8n0k%2FS%2Bw%2Fo0xBnE5BpU9Br0ZKo1ZLmFZOjEhesGljuzllqW50tH14aS14qm17mX9%2Bx4GAgUCEx02JySqOvpSXvI%2BYvp2orqmpzeGrQh%2Bsr6yssa2ttK6v0bKxMBy01bm4zLu5yry7yb29x77BzMPCxsLEzMXFxsXGx8fI3PLJ08vKysrKy8rL2s3MzczOH8LR0dHW19bX19fZ2dna2trc3Nzd3d3d3t3f39%2FgtZTg4ODi4uLj4%2BPlGxLl5eXm5ubnRzPn5%2Bfo6Ojp6enqfmzq6urr6%2Bvt7e3t7u3uDwvugwbu7u7v6Obv8fDz8%2FP09PT2igP29vb4%2BPj6y376%2Bu%2F7%2Bfv9%2Ff39%2Fv3%2BkAH%2FAwf%2FtwD%2F9wCyh1KfAAAAKXRSTlMABQ4VGykqLjVCTVNgdXuHj5Kaq62vt77ExNPX2%2Bju8vX6%2Bvr7%2FP7%2B%2FiiUMfUAAADTSURBVAjXBcFRTsIwHAfgX%2FtvOyjdYDUsRkFjTIwkPvjiOTyX9%2FAIJt7BF570BopEdHOOstHS%2BX0s439RGwnfuB5gSFOZAgDqjQOBivtGkCc7j%2B2e8XNzefWSu%2BsZUD1QfoTq0y6mZsUSvIkRoGYnHu6Yc63pDCjiSNE2kYLdCUAWVmK4zsxzO%2BQQFxNs5b479NHXopkbWX9U3PAwWAVSY%2FpZf1udQ7rfUpQ1CzurDPpwo16Ff2cMWjuFHX9qCV0Y0Ok4Jvh63IABUNnktl%2B6sgP%2BARIxSrT%2FMhLlAAAAAElFTkSuQmCC)](http://spiceprogram.org/oss-sponsorship)

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
