package main

import (
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"mama/internal/web"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	web.Launch()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {

	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}
	configurationDirectory := filepath.Dir(executable)
	configuration.Initialise(configurationDirectory)

	logwrapper.Initialise()

	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	mamasrv, error := service.New(prg, svcConfig)
	if error != nil {
		logwrapper.Log.Fatalf(error.Error())
	}
	logger, error = mamasrv.Logger(nil)
	if error != nil {
		logwrapper.Log.Fatalf(error.Error())
	}
	error = mamasrv.Run()
	if error != nil {
		logwrapper.Log.Fatalf(error.Error())
	}
}
