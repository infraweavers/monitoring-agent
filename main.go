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

func findConfigurationDirectory() string {
	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}

	executableFolder := filepath.Dir(executable)
	goSrcFolder := os.Getenv("GOPATH") + "\\src\\mama\\"

	_, error = os.Stat(executableFolder + "configuration.ini")
	if error == nil {
		return executableFolder
	}

	if os.IsNotExist(error) {
		_, error = os.Stat(goSrcFolder + "configuration.ini")
		if error == nil {
			return goSrcFolder
		}

		if os.IsNotExist(error) {
			statError := os.PathError{
				Op:   "stat",
				Path: executableFolder + "|" + goSrcFolder,
				Err:  error,
			}
			panic(statError)
		}
	}

	panic(error)
}

func main() {
	configuration.Initialise(findConfigurationDirectory())
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
