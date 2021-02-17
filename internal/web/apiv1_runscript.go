package web

import (
	"encoding/json"
	"fmt"
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net/http"
	"os/exec"
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

	processKiller := time.NewTimer(configuration.Settings.RequestTimeout)

	go func() {
		<-processKiller.C
		command.Process.Kill()
		logwrapper.Log.Warningf("Request Timed Out: '%s' %#v", scriptToRun.Path, scriptToRun.Args)
	}()

	output, error := command.CombinedOutput()

	if processKiller.Stop() == false {
		<-processKiller.C
	}

	if error != nil {
		if exitError, ok := error.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		}
	}

	return processResult(responseWriter, exitcode, string(output))
}

// APIV1RunscriptGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": ["executes a script sent via HTTP POST request"]}`))
}

// APIV1RunscriptPostHandler creates executes a script as specified in a http request and updates the http response with the result
func APIV1RunscriptPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	var script Script
	error := dec.Decode(&script)
	if error != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
