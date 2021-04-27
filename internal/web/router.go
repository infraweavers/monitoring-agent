package web

import (
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"

	// Blank import of pprof for side effect of loading its handlers
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

//NewRouter returns an HTTP multiplexor
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range declaredRoutes {
		logwrapper.LogDebugf("registering route Name: %s; Method: %s; Path: %s;", route.Name, route.Method, route.Pattern)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	logwrapper.LogInfof("configuration.Settings.LoadPprof: %t", configuration.Settings.LoadPprof)
	if configuration.Settings.LoadPprof {
		logwrapper.LogDebugf("registering '/debug/pprof' route due to configuration")
		router.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	}

	logwrapper.LogInfof("configuration.Settings.LogHTTPRequests: %t", configuration.Settings.LogHTTPRequests)
	if configuration.Settings.LogHTTPRequests {
		logwrapper.LogDebugf("appending httpRequestLogger middleware due to configuration")
		router.Use(httpRequestLogger)
	}

	logwrapper.LogInfof("configuration.Settings.LogHTTPResponses: %t", configuration.Settings.LogHTTPResponses)
	if configuration.Settings.LogHTTPResponses {
		logwrapper.LogDebugf("appending httpResponseLogger middleware due to configuration")
		router.Use(httpResponseLogger)
	}

	logwrapper.LogDebugf("appending IPFiltering middleware")
	router.Use(IPFiltering)

	logwrapper.LogDebugf("appending BasicAuth middleware")
	router.Use(BasicAuth)

	return router
}
