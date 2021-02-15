package main

import (
	"mama/internal/configuration"
	"mama/internal/httpserver"
	"mama/internal/logwrapper"
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
	logwrapper.Initialise()

	httpserver.Launch()
}
