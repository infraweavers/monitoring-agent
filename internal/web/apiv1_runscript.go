package web

import (
	"fmt"
	"net/http"
)

// APIV1RunscriptGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": ["executes a script sent via HTTP POST request"]}`))
}

// APIV1RunscriptPostHandler creates executes a script as specified in a http request and updates the http response with the result
func APIV1RunscriptPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := JsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if script.StdIn != "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint does not use stdin", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
