package web

import (
	"crypto/tls"
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net/http"
)

// Launch instantiates a multiplexer and uses it to configure and launch an HTTP server
func Launch() {

	tlsCert, loadError := tls.LoadX509KeyPair(configuration.Settings.CertificatePath, configuration.Settings.PrivateKeyPath)

	if loadError != nil {
		panic(loadError)
	}

	var requestTimeout = configuration.Settings.RequestTimeout

	router := NewRouter()
	handlerWithTimeout := http.TimeoutHandler(router, requestTimeout, "Request Timeout\n")

	server := &http.Server{
		Addr:    configuration.Settings.BindAddress,
		Handler: handlerWithTimeout,

		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
		WriteTimeout:      requestTimeout,
		ReadTimeout:       requestTimeout,
		ReadHeaderTimeout: requestTimeout,
		IdleTimeout:       requestTimeout,
	}

	logwrapper.Log.Info("Launching web server: https://" + configuration.Settings.BindAddress)
	logwrapper.Log.Info(server.ListenAndServeTLS("", ""))
}
