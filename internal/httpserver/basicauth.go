package basicauth

import (
	"mama/internal/configuration"
	"net/http"
)

func isKnownCredential(username string, password string) bool {
	return username == configuration.Settings.Username && password == configuration.Settings.Password
}

// IsAuthorised checks whether an http request is authorised. Returns a boolean
func IsAuthorised(responseWriter http.ResponseWriter, request *http.Request) bool {
	username, password, ok := request.BasicAuth()

	if !ok {
		responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
		responseWriter.WriteHeader(http.StatusUnauthorized)
		responseWriter.Write([]byte(`{"message": "Basic authentication required"}`))
		return false
	}

	if !isKnownCredential(username, password) {
		responseWriter.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
		responseWriter.WriteHeader(http.StatusForbidden)
		responseWriter.Write([]byte(`{"message": "Invalid username and/or password"}`))
		return false
	}

	return true
}
