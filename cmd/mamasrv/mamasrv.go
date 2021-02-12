package main

import (
	"mama/internal/httpserver"
	"os"
	"path/filepath"
)

func main() {
	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}
	configurationDirectory := filepath.Dir(executable)
	httpserver.Launch(configurationDirectory)
}
