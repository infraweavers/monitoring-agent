package apiv1

import (
    "net/http"
    "mama/internal/basicauth"
)

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    basicauth.IsAuthorised(w, r)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"endpoints": ["runscript", "info"]}`))
}

// w.WriteHeader(http.StatusMethodNotAllowed)
// w.Write([]byte( fmt.Sprintf( "{\"message\": \"HTTP %s not implemented\"}", r.Method ) ))