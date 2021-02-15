package httpserver

import (
	"net/http"
	"testing"
)

func TestDefaultHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("with no credentials, returns 401 unauthorized", func(t *testing.T) {
		output := TestHTTPRequest(t, BuildTestHTTPRequest(t, http.MethodGet, "/"))
		if output.ResponseStatus != http.StatusUnauthorized {
			t.Errorf("Test Failed: Expected: %d, Got: %d", http.StatusUnauthorized, output.ResponseStatus)
		}
	})

	t.Run("with incorrect credentials, returns 403 forbidden", func(t *testing.T) {
		output := TestHTTPRequestWithCredentials(t, BuildTestHTTPRequest(t, http.MethodGet, "/"), "bad_username", "bad_password")
		if output.ResponseStatus != http.StatusForbidden {
			t.Errorf("Test Failed: Expected: %d, Got: %d", http.StatusForbidden, output.ResponseStatus)
		}
	})
}
