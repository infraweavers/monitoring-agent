package main

import (
	"flag"
	"fmt"
	"log"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"monitoringagent/internal/web"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kardianos/service"
)

var maVersion string = "0.0.0"
var operatingSystem string = runtime.GOOS
var arch string = runtime.GOARCH
var goVersion = ""

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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var configDirectory string

	flag.StringVar(&configDirectory, "configurationDirectory", "", "Override the directory containing the configuration.")
	showVersion := flag.Bool("version", false, "Show version number and exit")
	flag.Parse()

	monitoringAgentVersionString := strings.Join([]string{"monitoring-agent " + maVersion, operatingSystem + " " + arch, goVersion}, "; ")

	if *showVersion {
		fmt.Printf("%s\n", monitoringAgentVersionString)
		os.Exit(0)
	}

	configuration.Initialise(configurationDirectory(configDirectory), monitoringAgentVersionString)

	logFile, _ := os.OpenFile(configuration.Settings.Logging.LogFilePath+".stdout", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	redirectStderr(logFile)

	logwrapper.Initialise(service.Interactive(), NewLine)

	serviceConfiguration := &service.Config{
		Name: "monitoring-agent",
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
	logwrapper.LogInfo("End of Main")
}
