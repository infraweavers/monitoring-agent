package basicauth

import (
	"net/http"
    "net/http/httptest"
    "testing"
)


var testCredentials = []struct {
    username    string
    password    string
    expected    bool
}{
    {"test",    "secret",       true},
    {"user",    "secret",       false},
    {"test",    "password",     false},
    {"user",    "password",     false},
    {"test",    "",             false},
    {"",        "password",     false},
    {"",        "",             false},
}

func TestIsAuthorised(t *testing.T) {

	request, _ := http.NewRequest("GET", "/", nil)
    responseWriter := httptest.NewRecorder()

    for _, test := range testCredentials {
        request.SetBasicAuth(test.username, test.password)

        if output := IsAuthorised(responseWriter, request); output != test.expected {
            t.Error("Test Failed: Input: {}:{}, Expected: {}, Got: {}", test.username, test.expected, output)
        }
    }
}