package httpserver

import (
	"net/http"
)

// APIV1InfoGetHandler creates a http response for the API /info path
func APIV1InfoGetHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json")
	reponseWriter.WriteHeader(http.StatusOK)
	reponseWriter.Write([]byte(`{"info": ["Not yet implemented. Will return agent info (version, memory consumption, cpu utilisation etc.)"]}`))
}
