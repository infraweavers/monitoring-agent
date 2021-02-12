package httpserver

import (
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// Launch instantiates a multiplexer and uses it to configure and launch an HTTP server
func Launch(configurationDirectory string) {

	certificatePath := filepath.FromSlash(configurationDirectory + "/server.crt")
	keyfilePath := filepath.FromSlash(configurationDirectory + "/server.key")

	tlsCert, loadError := tls.LoadX509KeyPair(certificatePath, keyfilePath)

	if loadError != nil {
		panic(loadError)
	}

	router := newRouter()

	server := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	log.Println("Launching web server: https://0.0.0.0:9000")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
