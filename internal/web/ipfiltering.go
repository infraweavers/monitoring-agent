package web

import (
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net"
	"net/http"
	"strings"
)

// IPFiltering is a HTTP handlefunc wrapper that enforces IP restrictions
func IPFiltering(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		remoteAddr := request.RemoteAddr
		remoteAddrComponents := strings.Split(remoteAddr, ":")

		remoteIp := net.ParseIP(remoteAddrComponents[0])
		allowedAddresses := configuration.Settings.AllowedAddresses

		allowAccess := false
		for x := 0; x < len(allowedAddresses); x++ {
			if allowedAddresses[x].Contains(remoteIp) {
				allowAccess = true
			}
		}

		if allowAccess {
			handler.ServeHTTP(responseWriter, request)
		} else {
			logwrapper.Log.Errorf("Blocked request from: %v", remoteAddr)
			responseWriter.WriteHeader(http.StatusForbidden)
			responseWriter.Write([]byte(`{"message": "Forbidden"}`))
		}
	})
}
