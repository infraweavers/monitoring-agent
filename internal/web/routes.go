package web

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
	route{"apiv1RunexecutableGet", "GET", "/v1/runexecutable", APIV1RunexecutableGetHandler},
	route{"apiv1RunexecutablePost", "POST", "/v1/runexecutable", APIV1RunexecutablePostHandler},
	route{"apiv1RunscriptstdinGet", "GET", "/v1/runscriptstdin", APIV1RunscriptstdinGetHandler},
	route{"apiv1RunscriptstdinPost", "POST", "/v1/runscriptstdin", APIV1RunscriptstdinPostHandler},
	route{"apiv1InfoGet", "GET", "/v1/info", APIV1InfoGetHandler},
	route{"apiv1Version", "GET", "/v1/version", APIV1VersionHandler},
}
