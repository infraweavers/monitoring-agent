package web

import (
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
)

func httpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		if configuration.Settings.LogHTTPRequests {
			logwrapper.LogHttpRequest(
				request.RemoteAddr,
				request.Host,
				request.Method,
				request.URL.String(),
				request.Header,
				request.Proto,
				request.ContentLength,
			)
		}

		handler.ServeHTTP(responseWriter, request)
	})
}
