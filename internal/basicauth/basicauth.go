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

func IsAuthorised(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()

	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Basic authentication required"}`))
		return false
	}

	if !isKnownCredential(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Access restricted"`)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message": "Invalid username and/or password"}`))
		return false
	}

	return true
}
