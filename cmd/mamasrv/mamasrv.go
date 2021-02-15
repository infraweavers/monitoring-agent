package main

import (
	"mama/internal/configuration"
	"mama/internal/httpserver"
	"mama/internal/logging"
	"os"
	"path/filepath"
)

func main() {
	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}

	configurationDirectory := filepath.Dir(executable)

	configuration.Initialise(configurationDirectory)
	logging.Initialise()

	httpserver.Launch()
}
