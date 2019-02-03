# JWT authentication proxy

**UNDER CONSTRUCTION, NOT READY TO BE USED YET**

[![CircleCI](https://circleci.com/gh/janivihervas/authproxy.svg?style=svg)](https://circleci.com/gh/janivihervas/authproxy)
[![codecov](https://codecov.io/gh/janivihervas/authproxy/branch/master/graph/badge.svg)](https://codecov.io/gh/janivihervas/authproxy)

[![Go Report Card](https://goreportcard.com/badge/github.com/janivihervas/authproxy)](https://goreportcard.com/report/github.com/janivihervas/authproxy)
[![GoDoc](https://godoc.org/github.com/janivihervas/authproxy?status.svg)](https://godoc.org/github.com/janivihervas/authproxy)
[![Sponsored](https://img.shields.io/badge/chilicorn-sponsored-brightgreen.svg?logo=data%3Aimage%2Fpng%3Bbase64%2CiVBORw0KGgoAAAANSUhEUgAAAA4AAAAPCAMAAADjyg5GAAABqlBMVEUAAAAzmTM3pEn%2FSTGhVSY4ZD43STdOXk5lSGAyhz41iz8xkz2HUCWFFhTFFRUzZDvbIB00Zzoyfj9zlHY0ZzmMfY0ydT0zjj92l3qjeR3dNSkoZp4ykEAzjT8ylUBlgj0yiT0ymECkwKjWqAyjuqcghpUykD%2BUQCKoQyAHb%2BgylkAyl0EynkEzmkA0mUA3mj86oUg7oUo8n0k%2FS%2Bw%2Fo0xBnE5BpU9Br0ZKo1ZLmFZOjEhesGljuzllqW50tH14aS14qm17mX9%2Bx4GAgUCEx02JySqOvpSXvI%2BYvp2orqmpzeGrQh%2Bsr6yssa2ttK6v0bKxMBy01bm4zLu5yry7yb29x77BzMPCxsLEzMXFxsXGx8fI3PLJ08vKysrKy8rL2s3MzczOH8LR0dHW19bX19fZ2dna2trc3Nzd3d3d3t3f39%2FgtZTg4ODi4uLj4%2BPlGxLl5eXm5ubnRzPn5%2Bfo6Ojp6enqfmzq6urr6%2Bvt7e3t7u3uDwvugwbu7u7v6Obv8fDz8%2FP09PT2igP29vb4%2BPj6y376%2Bu%2F7%2Bfv9%2Ff39%2Fv3%2BkAH%2FAwf%2FtwD%2F9wCyh1KfAAAAKXRSTlMABQ4VGykqLjVCTVNgdXuHj5Kaq62vt77ExNPX2%2Bju8vX6%2Bvr7%2FP7%2B%2FiiUMfUAAADTSURBVAjXBcFRTsIwHAfgX%2FtvOyjdYDUsRkFjTIwkPvjiOTyX9%2FAIJt7BF570BopEdHOOstHS%2BX0s439RGwnfuB5gSFOZAgDqjQOBivtGkCc7j%2B2e8XNzefWSu%2BsZUD1QfoTq0y6mZsUSvIkRoGYnHu6Yc63pDCjiSNE2kYLdCUAWVmK4zsxzO%2BQQFxNs5b479NHXopkbWX9U3PAwWAVSY%2FpZf1udQ7rfUpQ1CzurDPpwo16Ff2cMWjuFHX9qCV0Y0Ok4Jvh63IABUNnktl%2B6sgP%2BARIxSrT%2FMhLlAAAAAElFTkSuQmCC)](http://spiceprogram.org/oss-sponsorship)

<!-- toc -->

<!-- tocstop -->

Use cases:

- As a standalone application (reverse proxy) for authenticating requests to your application
- As a Go middleware for authenticating requests to your server

Features:

- Validating JWT access tokens
- Refreshing JWT access tokens with refresh tokens

Works the same way as [bitly/oauth2_proxy](https://github.com/bitly/oauth2_proxy) but with a few additions:

- Refreshes access tokens.
- Haves an option of **not** redirecting every request to OIDC/OAuth2 provider's login page. E.g. redirecting API calls with expired authentication from frontend login page doesn't work. Requests not redirected to the login page will receive a response with 401 or 403 as status code and a JSON with `redirectUrl` field so the frontend can redirect the browser itself.
- Redirects back to original url once authentication is renewed

See [flow diagram](doc/flow-diagram.md) for details.
