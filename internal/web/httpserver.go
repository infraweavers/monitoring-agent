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

	logwrapper.LogInfof("configuration.Settings.HTTPRequestTimeout: %v", configuration.Settings.HTTPRequestTimeout)
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

	logwrapper.LogInfof("configuration.Settings.DisableHTTPs: %t", configuration.Settings.DisableHTTPs)
	if configuration.Settings.DisableHTTPs {
		logwrapper.LogCriticalf("!! The HTTP server is running insecurely due to 'configuration.Settings.DisableHTTPs'='%t'. This is not a recommended setting !!", configuration.Settings.DisableHTTPs)
		logwrapper.LogCritical("!! Re-enable HTTPs by setting 'DisableHTTPs' to 'true' as soon as possible !!")
		err := server.ListenAndServe()
		if err != nil {
			logwrapper.Log.Panicf(err.Error())
		}
		return
	}

	err := server.ListenAndServeTLS("", "")
	if err != nil {
		logwrapper.Log.Panicf(err.Error())
	}
}
