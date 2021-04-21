package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("with no credentials, returns 401 unauthorized", func(t *testing.T) {
		output := TestHTTPRequest(t, BuildTestHTTPRequest(t, http.MethodGet, "/"))
		assert.Equal(t, http.StatusUnauthorized, output.ResponseStatus)
	})

	t.Run("with incorrect credentials, returns 403 forbidden", func(t *testing.T) {
		output := TestHTTPRequestWithCredentials(t, BuildTestHTTPRequest(t, http.MethodGet, "/"), "bad_username", "bad_password")
		assert.Equal(t, http.StatusForbidden, output.ResponseStatus)
	})
}
