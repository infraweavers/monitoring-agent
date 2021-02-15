package apiv1

import (
	"mama/internal/basicauth"
	"net/http"
)

// InfoGetHandler creates a http response for the API /info path
func InfoGetHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(reponseWriter, request)
	reponseWriter.WriteHeader(http.StatusOK)
	reponseWriter.Write([]byte(`{"info": ["Not yet implemented. Will return agent info (version, memory consumption, cpu utilisation etc.)"]}`))
}
