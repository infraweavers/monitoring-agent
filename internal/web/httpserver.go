package web

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
)

// LaunchServer instantiates a multiplexer and uses it to configure and launch an HTTP server
func LaunchServer() {

	tlsCert, loadError := tls.LoadX509KeyPair(configuration.Settings.CertificatePath, configuration.Settings.PrivateKeyPath)

	if loadError != nil {
		logwrapper.Log.Panicf(loadError.Error())
	}

	var requestTimeout = configuration.Settings.HTTPRequestTimeout

	router := NewRouter()
	handlerWithTimeout := http.TimeoutHandler(router, requestTimeout, "Request Timeout\n")

	var clientAuthenticationMethod = tls.NoClientCert

	clientCACertPool := x509.NewCertPool()
	if configuration.Settings.UseClientCertificates {
		clientAuthenticationMethod = tls.RequireAndVerifyClientCert
		certificateContent, clientCALoaderror := ioutil.ReadFile(configuration.Settings.ClientCertificateCAFile)

		if clientCALoaderror != nil {
			logwrapper.LogCriticalf("Unable to read ClientCertificateCAFile from %s ", configuration.Settings.ClientCertificateCAFile)
			panic(clientCALoaderror)
		}

		clientCACertPool.AppendCertsFromPEM(certificateContent)
	}

	server := &http.Server{
		Addr:    configuration.Settings.BindAddress,
		Handler: handlerWithTimeout,

		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			ClientAuth:   clientAuthenticationMethod,
			ClientCAs:    clientCACertPool,
		},
		WriteTimeout:      requestTimeout,
		ReadTimeout:       requestTimeout,
		ReadHeaderTimeout: requestTimeout,
		IdleTimeout:       requestTimeout,
	}

	logwrapper.LogInfof("Launching web server: https://%s", configuration.Settings.BindAddress)
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		logwrapper.Log.Panicf(err.Error())
	}
}
