package apiv1

import (
	"mama/internal/basicauth"
	"net/http"
)

// HomeGetHandler creates a http response for the API root path
func HomeGetHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(reponseWriter, request)
	reponseWriter.WriteHeader(http.StatusOK)
	reponseWriter.Write([]byte(`{"endpoints": ["runscript", "info"]}`))
}

// w.WriteHeader(http.StatusMethodNotAllowed)
// w.Write([]byte( fmt.Sprintf( "{\"message\": \"HTTP %s not implemented\"}", r.Method ) ))
