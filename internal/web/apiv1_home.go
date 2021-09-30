package web

import (
	"net/http"
)

// APIV1HomeGetHandler creates a http response for the API root path
func APIV1HomeGetHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json")
	reponseWriter.WriteHeader(http.StatusOK)
	reponseWriter.Write([]byte(`{"endpoints": ["runexecutable", "runscriptstdin", "info"]}`))
}
