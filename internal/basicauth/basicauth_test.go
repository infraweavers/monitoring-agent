package basicauth

import (
    "testing"
)


func TestIsKnownCredential(t *testing.T) {
    var tests = []struct {
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

    for _, test := range tests {
        if output := IsKnownCredential(test.username, test.password); output != test.expected {
            t.Error("Test Failed: Input: {}:{}, Expected: {}, Got: {}", test.username, test.expected, output)
        }
    }
}