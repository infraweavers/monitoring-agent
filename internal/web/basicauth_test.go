package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAuthorised(t *testing.T) {
	var testCases = []struct {
		credential TestCredential
		expected   int
	}{
		{TestCredential{Username: "test", Password: "secret"}, http.StatusOK},
		{TestCredential{Username: "user", Password: "secret"}, http.StatusForbidden},
		{TestCredential{Username: "test", Password: "password"}, http.StatusForbidden},
		{TestCredential{Username: "user", Password: "password"}, http.StatusForbidden},
		{TestCredential{Username: "test", Password: ""}, http.StatusForbidden},
		{TestCredential{Username: "", Password: "password"}, http.StatusForbidden},
		{TestCredential{Username: "", Password: ""}, http.StatusForbidden},
	}

	TestSetup()
	defer TestTeardown()

	for _, testCase := range testCases {
		t.Run("returns appropriate responses with valid/invalid credentials", func(t *testing.T) {
			output := TestHTTPRequestWithCredentials(t, BuildTestHTTPRequest(t, http.MethodGet, "/v1"), testCase.credential.Username, testCase.credential.Password)
			assert.Equal(t, output.ResponseStatus, testCase.expected)
		})
	}

	t.Run("returns false when no credentials are supplied", func(t *testing.T) {
		expected := http.StatusUnauthorized
		output := TestHTTPRequest(t, BuildTestHTTPRequest(t, http.MethodGet, "/v1"))
		assert.Equal(t, output.ResponseStatus, expected)
	})
}
