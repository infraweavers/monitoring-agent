package httpserver

import (
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const requestTimeout = 30 * time.Second

// Launch instantiates a multiplexer and uses it to configure and launch an HTTP server
func Launch(configurationDirectory string) {

	certificatePath := filepath.FromSlash(configurationDirectory + "/server.crt")
	keyfilePath := filepath.FromSlash(configurationDirectory + "/server.key")

	tlsCert, loadError := tls.LoadX509KeyPair(certificatePath, keyfilePath)

	if loadError != nil {
		panic(loadError)
	}

	router := NewRouter()
	handlerWithTimeout := http.TimeoutHandler(router, requestTimeout, "Request Timeout\n")

	server := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: handlerWithTimeout,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
		WriteTimeout:      requestTimeout,
		ReadTimeout:       requestTimeout,
		ReadHeaderTimeout: requestTimeout,
		IdleTimeout:       requestTimeout,
	}
	log.Println("Launching web server: https://0.0.0.0:9000")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
