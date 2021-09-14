package main

import (
	"flag"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"monitoringagent/internal/web"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

type program struct{}

func (program *program) Start(s service.Service) error {
	logwrapper.LogInfo("Service Starting")
	go program.run()
	return nil
}

func (program *program) run() {
	logwrapper.LogInfo("calling web.LaunchServer")
	web.LaunchServer()
}

func (program *program) Stop(s service.Service) error {
	web.KillAllRunningProcs()
	logwrapper.LogInfo("Service Stopping")
	return nil
}

func configurationDirectory(commandLineDirectory string) string {

	if commandLineDirectory != "" {
		return commandLineDirectory
	}

	executable, error := os.Executable()
	if error != nil {
		panic(error)
	}

	executableFolder := filepath.Dir(executable)
	_, error = os.Stat(filepath.FromSlash(executableFolder + "/configuration.json"))
	if error == nil {
		return executableFolder
	}

	if os.IsNotExist(error) {

		workingDirectory, _ := os.Getwd()
		_, error = os.Stat(filepath.FromSlash(workingDirectory + "/configuration.json"))
		if error == nil {
			return workingDirectory
		}

		if os.IsNotExist(error) {
			statError := os.PathError{
				Op:   "stat",
				Path: filepath.FromSlash(workingDirectory + "/configuration.json"),
				Err:  error,
			}
			panic(statError)
		}
	}

	panic(error)
}

func main() {

	var configDirectory string

	flag.StringVar(&configDirectory, "configurationDirectory", "", "Override the directory containing the configuration.")
	flag.Parse()

	configuration.Initialise(configurationDirectory(configDirectory))
	logwrapper.Initialise(service.Interactive(), NewLine)

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
