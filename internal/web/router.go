package web

import (
	"mama/internal/configuration"
	"net/http"

	// Blank import of pprof for side effect of loading its handlers
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

// MiddlewareFunc is an interface required for middleware wrapped around HTTP handler functions
type MiddlewareFunc func(http.Handler) http.Handler

//NewRouter returns an HTTP multiplexor
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range declaredRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	if configuration.Settings.LoadPprof {
		router.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	}
	router.Use(BasicAuth)
	return router
}
