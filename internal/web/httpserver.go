package web

import (
	"crypto/tls"
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net/http"
)

// LaunchServer instantiates a multiplexer and uses it to configure and launch an HTTP server
func LaunchServer() {

	tlsCert, loadError := tls.LoadX509KeyPair(configuration.Settings.CertificatePath, configuration.Settings.PrivateKeyPath)

	if loadError != nil {
		panic(loadError)
	}

	var requestTimeout = configuration.Settings.RequestTimeout

	router := NewRouter()
	handlerWithTimeout := http.TimeoutHandler(router, requestTimeout, "Request Timeout\n")

	var clientAuthenticationMethod = tls.NoClientCert
	if configuration.Settings.UseClientCertificates {
		clientAuthenticationMethod = tls.RequireAndVerifyClientCert
	}

	server := &http.Server{
		Addr:    configuration.Settings.BindAddress,
		Handler: handlerWithTimeout,

		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			ClientAuth:   clientAuthenticationMethod,
		},
		WriteTimeout:      requestTimeout,
		ReadTimeout:       requestTimeout,
		ReadHeaderTimeout: requestTimeout,
		IdleTimeout:       requestTimeout,
	}

	logwrapper.Log.Info("Launching web server: https://" + configuration.Settings.BindAddress)
	logwrapper.Log.Info(server.ListenAndServeTLS("", ""))
}
