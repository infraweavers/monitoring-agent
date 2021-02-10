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

func TestIsKnownCredential(t *testing.T) {
    for _, test := range testCredentials {
        if output := IsKnownCredential(test.username, test.password); output != test.expected {
            t.Error("Test Failed: Input: {}:{}, Expected: {}, Got: {}", test.username, test.expected, output)
        }
    }
}

func TestIsAuthorised(t *testing.T) {

	request, _ := http.NewRequest("GET", "/", nil)
    wwriter := http.ResponseWriter

    for _, test := range testCredentials {
        req.SetBasicAuth(username, password)

        if output := IsAuthorised(writer, request)); output != test.expected {
            t.Error("Test Failed: Input: {}:{}, Expected: {}, Got: {}", test.username, test.expected, output)
        }
    }
}