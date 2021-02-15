package web

import (
	"net/http"
)

// DefaultHandler generates the http response for that application root path ("/")
func DefaultHandler(responseWriter http.ResponseWriter, response *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": "/v1/"}`))

}
