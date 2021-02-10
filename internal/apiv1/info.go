package apiv1

import (
    "net/http"
    "mama/internal/basicauth"
)

func InfoGetHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    basicauth.IsAuthorised(w, r)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"endpoints": ["Not yet implemented. Will return agent info and status?"]}`))
}