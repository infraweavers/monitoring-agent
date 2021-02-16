package main

import (
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"mama/internal/web"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

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
	_, error = os.Stat(filepath.FromSlash(executableFolder + "/configuration.ini"))
	if error == nil {
		return executableFolder
	}

	if os.IsNotExist(error) {
		goSrcFolder := filepath.FromSlash(os.Getenv("GOPATH") + "/src/mama/")
		_, error = os.Stat(filepath.FromSlash(goSrcFolder + "/configuration.ini"))
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
	logwrapper.Initialise(service.Interactive())

	serviceConfiguration := &service.Config{
		Name: "Monitoring Agent",
	}

	program := &program{}

	serverInstance, serverError := service.New(program, serviceConfiguration)
	if serverError != nil {
		logwrapper.Log.Fatalf(serverError.Error())
	}

	instanceError := serverInstance.Run()
	if instanceError != nil {
		logwrapper.Log.Fatalf(instanceError.Error())
	}
}
