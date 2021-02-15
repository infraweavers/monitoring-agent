package web

import (
	"crypto/subtle"
	"mama/internal/configuration"
	"net/http"
)

func isKnownValidCredential(suppliedUsername string, suppliedPassword string) bool {
	if subtle.ConstantTimeCompare([]byte(configuration.Settings.Username), []byte(suppliedUsername)) != 1 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(configuration.Settings.Password), []byte(suppliedPassword)) != 1 {
		return false
	}
	return true
}

// BasicAuth is a HTTP handlefunc wrapper that handles HTTP basic authorization
func BasicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		username, password, ok := request.BasicAuth()

		if !ok {
			responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
			responseWriter.WriteHeader(http.StatusUnauthorized)
			responseWriter.Write([]byte(`{"message": "Basic authentication required"}`))
			return
		}

		if !isKnownValidCredential(username, password) {
			responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
			responseWriter.WriteHeader(http.StatusForbidden)
			responseWriter.Write([]byte(`{"message": "Invalid username and/or password"}`))
			return
		}

		handler.ServeHTTP(responseWriter, request)
	})
}
