package basicauth

import (
	"net/http"
)

func isKnownCredential(username string, password string) bool {
	var users = map[string]string{
		"test": "secret",
	}

	_password, ok := users[username]
	if !ok {
		return false
	}

	return password == _password
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
