package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type credential struct {
	username string
	password string
}

func tester(t *testing.T, cred credential) int {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	responseWriter := httptest.NewRecorder()

	if cred.username != "" && cred.password != "" {
		request.SetBasicAuth("cred.username", "cred.password")
	}

	DefaultHandler(responseWriter, request)
	return responseWriter.Result().StatusCode
}

func TestDefaultHandler(t *testing.T) {
	var testCases = []struct {
		name           string
		testCredential credential
		expectedResult int
	}{
		{"with no credentials, returns 401 unauthorized", credential{"", ""}, http.StatusUnauthorized},
		{"with incorrect credentials, returns 401 forbidden", credential{"foo", "bah"}, http.StatusForbidden},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			outputStatusCode := tester(t, testCase.testCredential)

			if outputStatusCode != testCase.expectedResult {
				t.Error("Test Failed: Expected: {}, Got: {}", testCase.expectedResult, outputStatusCode)
			}
		})
	}
}
