package httpserver

import (
	"mama/internal/apiv1"
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
	route{"Index", "GET", "/", DefaultHandler},
	// API Version 1
	route{"apiv1Index", "GET", "/v1", apiv1.HomeGetHandler},
	route{"apiv1Runscript", "GET", "/v1/runscript", apiv1.RunscriptGetHandler},
	route{"apiv1Runscript", "POST", "/v1/runscript", apiv1.RunscriptPostHandler},
	route{"apiv1Info", "GET", "/v1/info", apiv1.InfoGetHandler},
}
