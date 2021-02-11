package httpserver

import (
	"mama/internal/basicauth"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(w, r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"endpoints": "/v1/"}`))

}
