package httpserver

import (
	"mama/internal/basicauth"
	"net/http"
)

// DefaultHandler generates the http response for that application route path ("/")
func DefaultHandler(responseWriter http.ResponseWriter, response *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(responseWriter, response)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": "/v1/"}`))

}
