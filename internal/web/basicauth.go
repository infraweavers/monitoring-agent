package web

import (
	"crypto/subtle"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
)

func isKnownValidCredential(suppliedUsername string, suppliedPassword string) bool {
	if subtle.ConstantTimeCompare([]byte(configuration.Settings.Authentication.Username), []byte(suppliedUsername)) != 1 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(configuration.Settings.Authentication.Password), []byte(suppliedPassword)) != 1 {
		return false
	}
	return true
}

// BasicAuth is a HTTP handlefunc wrapper that handles HTTP basic authorization
func BasicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		username, password, ok := request.BasicAuth()

		if !ok {
			logwrapper.LogDebugf("Request received without Authorization header from: %v", request.RemoteAddr)
			responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
			responseWriter.WriteHeader(http.StatusUnauthorized)
			responseWriter.Write([]byte(`{"message": "Basic authentication required"}`))
			return
		}

		if !isKnownValidCredential(username, password) {
			logwrapper.LogInfof("Request received invalid basic auth from: %v", request.RemoteAddr)
			responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
			responseWriter.WriteHeader(http.StatusForbidden)
			responseWriter.Write([]byte(`{"message": "Invalid username and/or password"}`))
			return
		}

		handler.ServeHTTP(responseWriter, request)
	})
}
