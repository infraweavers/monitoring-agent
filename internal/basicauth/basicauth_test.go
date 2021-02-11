package basicauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type credential struct {
	username string
	password string
}

func tester(t *testing.T, cred credential) bool {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	responseWriter := httptest.NewRecorder()
	request.SetBasicAuth(cred.username, cred.password)

	return IsAuthorised(responseWriter, request)
}

func TestIsAuthorised(t *testing.T) {
	var testCases = []struct {
		testCredential credential
		expected       bool
	}{
		{credential{"test", "secret"}, true},
		{credential{"user", "secret"}, false},
		{credential{"test", "password"}, false},
		{credential{"user", "password"}, false},
		{credential{"test", ""}, false},
		{credential{"", "password"}, false},
		{credential{"", ""}, false},
	}

	for _, testCase := range testCases {
		t.Run("returns appropriate responses with valid/invalid credentials", func(t *testing.T) {
			output := tester(t, testCase.testCredential)

			if output != testCase.expected {
				t.Error("Test Failed: Input: {}, Expected: {}, Got: {}", testCase.testCredential, testCase.expected, output)
			}
		})
	}

	t.Run("returns false when no credentials are supplied", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseWriter := httptest.NewRecorder()
		expected := false
		output := IsAuthorised(responseWriter, request)

		if output != expected {
			t.Error("Test Failed: Expected: {}, Got: {}", expected, output)
		}
	})
}
