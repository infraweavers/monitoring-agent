package httpserver

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var declaredRoutes = routes{
	// Default
	route{"IndexGet", "GET", "/", DefaultHandler},
	// API Version 1
	route{"apiv1IndexGet", "GET", "/v1", APIV1HomeGetHandler},
	route{"apiv1RunscriptGet", "GET", "/v1/runscript", APIV1RunscriptGetHandler},
	route{"apiv1RunscriptPost", "POST", "/v1/runscript", APIV1RunscriptPostHandler},
	route{"apiv1InfoGet", "GET", "/v1/info", APIV1InfoGetHandler},
}
