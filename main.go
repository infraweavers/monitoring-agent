package main

import (
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"mama/internal/web"
	"os"
	"path/filepath"
	"runtime"

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
	directorySeparator := "/"
	if runtime.GOOS == "windows" {
		directorySeparator = "\\"
	}

	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}

	executableFolder := filepath.Dir(executable)
	goSrcFolder := os.Getenv("GOPATH") + directorySeparator + "src" + directorySeparator + "mama" + directorySeparator

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
		Name:        "Monitoring Agent",
		DisplayName: "Maintainable Monitoring Agent",
		Description: "Cross platform monitoring agent written in Go",
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
