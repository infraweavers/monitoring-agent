package httpserver

import (
    "net/http"
    "mama/internal/basicauth"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    basicauth.IsAuthorised(w, r)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"endpoints": "/v1/"}`))

}