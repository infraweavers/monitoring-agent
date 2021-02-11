package basicauth

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestIsAuthorised(t *testing.T) {
    t.Run("returns appropriate responses with valid/invalid credentials", func(t *testing.T) {
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
        
        request, _ := http.NewRequest("GET", "/", nil)
        responseWriter := httptest.NewRecorder()

        for _, test := range testCredentials {
            request.SetBasicAuth(test.username, test.password)
            output := IsAuthorised(responseWriter, request)

            if output != test.expected {
                t.Error("Test Failed: Input: {}:{}, Expected: {}, Got: {}", test.username, test.expected, output)
            }
        }
    })
    
    t.Run("returns false when no credentials are supplied", func(t *testing.T) {
        request, _ := http.NewRequest("GET", "/", nil)
        responseWriter := httptest.NewRecorder()
        
        expected := false
        output := IsAuthorised(responseWriter, request)
        
        if output != expected {
            t.Error("Test Failed: Expected: {}, Got: {}", expected, output)
        }
    })
}