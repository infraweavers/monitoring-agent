package web

import (
	"encoding/json"
	"fmt"
	"mama/internal/configuration"
	"mama/internal/logwrapper"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/openpgp"
)

// Script represents an object submitted to the runscript endpoint
type Script struct {
	Path      string   `json:"path"`
	Args      []string `json:"args"`
	StdIn     string   `json:"stdin"`
	Signature string   `json:"signature"`
}

// Result represents the object returned from the runscript endpoint
type Result struct {
	Exitcode int    `json:"exitcode"`
	Output   string `json:"output"`
}

type endpointDescription struct {
	Endpoint        string `json:"endpoint"`
	Description     string `json:"description"`
	MandatoryFields string `json:"mandatoryfields"`
	OptionalFields  string `json:"optionalfields"`
	ExampleRequest  string `json:"exampleRequest"`
	ExampleResponse string `json:"exampleResponse"`
}

type safeCollection struct {
	collection map[*exec.Cmd]bool
	mutex      sync.Mutex
}

var runningProcesses = safeCollection{}

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

func jsonDecodeScript(responseWriter http.ResponseWriter, request *http.Request) (Script, error) {
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	var script Script
	error := dec.Decode(&script)
	if error != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request", http.StatusBadRequest)))
		return script, error
	}

	return script, nil
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

	if scriptToRun.StdIn != "" {
		command.Stdin = strings.NewReader(scriptToRun.StdIn)
	}

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

func verifySignature(stdin string, signature string) bool {

	stdInReader := strings.NewReader(stdin)
	signatureReader := strings.NewReader(signature)
	signer, error := openpgp.CheckDetachedSignature(configuration.Settings.PublicKeyRing, stdInReader, signatureReader)
	if signer != nil && error != nil {
		return true
	}
	return false
}
