package web

import (
	"monitoringagent/internal/configuration"
	"net/http"

	// Blank import of pprof for side effect of loading its handlers
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

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
	router.Use(IPFiltering)
	router.Use(BasicAuth)

	if configuration.Settings.LogHTTPRequests {
		router.Use(httpRequestLogger)
	}

	if configuration.Settings.LogHTTPResponses {
		router.Use(httpResponseLogger)
	}

	return router
}
