package web

import (
	"monitoringagent/internal/logwrapper"
	"net/http"
)

// IPFiltering is a HTTP handlefunc wrapper that enforces IP restrictions
func IPFiltering(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if verifyRemoteHost(request.RemoteAddr) {
			logwrapper.LogDebugf("Allowed request due to IP restrictions from: %v", request.RemoteAddr)
			handler.ServeHTTP(responseWriter, request)
		} else {
			logwrapper.LogErrorf("Blocked request due to IP restrictions from: %v", request.RemoteAddr)
			http.Error(responseWriter, `Forbidden`, http.StatusForbidden)
		}
	})
}
