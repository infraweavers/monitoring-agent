package web

import (
	"encoding/json"
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

// Script represents an object submitted to the runscript endpoint
type Script struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

// Result represents the object returned from the runscript endpoint
type Result struct {
	Exitcode int    `json:"exitcode"`
	Output   string `json:"output"`
}

type SafeCollection struct {
	collection map[*exec.Cmd]bool
	mutex      sync.Mutex
}

var runningProcesses = SafeCollection{}

func init() {
	runningProcesses.collection = make(map[*exec.Cmd]bool)
}

// KillAllRunningProcs self explanatory
func KillAllRunningProcs() {
	logwrapper.Log.Info("Killing all procs")
	runningProcesses.mutex.Lock()
	for item := range runningProcesses.collection {
		logwrapper.Log.Debugf("Killing: %#v", item)
		item.Process.Kill()
	}
	runningProcesses.mutex.Unlock()
}

func processResult(responseWriter http.ResponseWriter, exitCode int, output string) []byte {

	result := Result{
		Exitcode: exitCode,
		Output:   output,
	}

	resultJSON, error := json.Marshal(result)
	if error != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return processResult(responseWriter, 3, error.Error())
	}

	return resultJSON
}

func runScript(responseWriter http.ResponseWriter, scriptToRun Script) []byte {

	var exitcode int = 0

	command := exec.Command(scriptToRun.Path, scriptToRun.Args...)

	runningProcesses.mutex.Lock()
	runningProcesses.collection[command] = true
	runningProcesses.mutex.Unlock()

	processKiller := time.AfterFunc(configuration.Settings.RequestTimeout, func() {
		command.Process.Kill()
		logwrapper.Log.Warningf("Request Timed Out: '%s' %#v", scriptToRun.Path, scriptToRun.Args)
	})

	output, error := command.CombinedOutput()
	processKiller.Stop()

	runningProcesses.mutex.Lock()
	delete(runningProcesses.collection, command)
	runningProcesses.mutex.Unlock()

	if error != nil {
		if exitError, ok := error.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		}
	}

	return processResult(responseWriter, exitcode, string(output))
}
