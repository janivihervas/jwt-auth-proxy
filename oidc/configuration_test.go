package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration_Valid(t *testing.T) {
	var (
		cupaloy = cupaloy.New(
			cupaloy.CreateNewAutomatically(true),
			cupaloy.FailOnUpdate(false),
			cupaloy.ShouldUpdate(func() bool {
				return true
			}),
		)
		config Configuration
		err    error
	)

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("Issuer", err.Error())
	assert.NoError(t, err)

	config.Issuer = "foo"

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("AuthorizationEndpoint", err.Error())
	assert.NoError(t, err)

	config.AuthorizationEndpoint = "foo"

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("JWKSURI", err.Error())
	assert.NoError(t, err)

	config.JWKSURI = "foo"

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("ResponseTypesSupported", err.Error())
	assert.NoError(t, err)

	config.ResponseTypesSupported = []string{"foo"}

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("SubjectTypesSupported", err.Error())
	assert.NoError(t, err)

	config.SubjectTypesSupported = []string{"foo"}

	err = config.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("IDTokenSigningAlgValuesSupported", err.Error())
	assert.NoError(t, err)

	config.IDTokenSigningAlgValuesSupported = []string{"foo"}

	err = config.Valid()
	assert.NoError(t, err)
}

func TestConfiguration_FillDefaultValuesIfEmpty(t *testing.T) {
	var config Configuration
	config.FillDefaultValuesIfEmpty()

	b := true
	assert.Equal(t, Configuration{
		ResponseModesSupported:            []string{"query", "fragment"},
		GrantTypesSupported:               []string{"authorization_code", "implicit"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"},
		RequestURIParameterSupported:      &b,
	}, config)
}

func TestGetOpenIDConnectConfiguration(t *testing.T) {
	var (
		ctx       = context.Background()
		server500 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		config = Configuration{
			Issuer: "iss",
		}
		serverBadResponse = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{`))
			assert.NoError(t, err)
		}))
		serverOK = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(config)
			assert.NoError(t, err)
		}))
	)
	defer server500.Close()
	defer serverBadResponse.Close()
	defer serverOK.Close()

	_, err := GetOpenIDConnectConfiguration(ctx, nil, "http://192.168.0.%31/")
	assert.Error(t, err)

	_, err = GetOpenIDConnectConfiguration(ctx, nil, "")
	assert.Error(t, err)

	_, err = GetOpenIDConnectConfiguration(ctx, nil, server500.URL)
	assert.Error(t, err)

	_, err = GetOpenIDConnectConfiguration(ctx, nil, serverBadResponse.URL)
	assert.Error(t, err)

	newConfig, err := GetOpenIDConnectConfiguration(ctx, nil, serverOK.URL)
	assert.NoError(t, err)
	assert.NotEqual(t, config, newConfig)

	config.FillDefaultValuesIfEmpty()
	assert.Equal(t, config, newConfig)
}
