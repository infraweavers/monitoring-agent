package httpserver
 
import (
    "net/http"
 	"internal/apiv1"
)
 
type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}
 
type Routes []Route

var routes = Routes{
	// Default
    Route{ "Index",				"GET",	"/",				defaultHandler, 			},
	// API Version 1
    Route{ "apiv1Index",		"GET",  "/v1",         		apiv1.HomeGetHandler,		},
    Route{ "apiv1Runscript",    "GET",  "/v1/runscript",    apiv1.RunscriptGetHandler,  },
	Route{ "apiv1Runscript",    "POST", "/v1/runscript",    apiv1.RunscriptPostHandler, },
	Route{ "apiv1Info",         "GET",  "/v1/info",         apiv1.InfoGetHandler,		},
}

