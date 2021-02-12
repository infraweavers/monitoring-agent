package apiv1

import (
	"encoding/json"
	"fmt"
	"mama/internal/basicauth"
	"net/http"
	"os/exec"
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

	resultJSON, err := json.Marshal(result)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return processResult(responseWriter, 3, err.Error())
	}

	return resultJSON
}

func runScript(responseWriter http.ResponseWriter, scriptToRun Script) []byte {
	var exitcode int = 0
	cmd := exec.Command(scriptToRun.Path, scriptToRun.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		}
	}

	return processResult(responseWriter, exitcode, string(output))
}

// RunscriptGetHandler creates a http response for the API /runscript http get requests
func RunscriptGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(responseWriter, request)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": ["executes a script sent via HTTP POST request"]}`))
}

// RunscriptPostHandler creates executes a script as specified in a http request and updates the http response with the result
func RunscriptPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	basicauth.IsAuthorised(responseWriter, request)

	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	var script Script
	err := dec.Decode(&script)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
