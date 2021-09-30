package web

import (
	"encoding/json"
	"fmt"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jedisct1/go-minisign"
)

// Script represents an object submitted to the runexecutable endpoint
type Script struct {
	Path           string   `json:"path"`
	Args           []string `json:"args,omitempty"`
	StdIn          string   `json:"stdin,omitempty"`
	StdInSignature string   `json:"stdinsignature,omitempty"`
	Timeout        string   `json:"timeout,omitempty"`
}

// Result represents the object returned from the runexecutable endpoint
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
	logwrapper.LogInfo("Killing all procs")
	runningProcesses.mutex.Lock()
	for item := range runningProcesses.collection {
		logwrapper.LogDebugf("Killing: %#v", item)
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
		logwrapper.LogWarningf("Failed JSON decode: '%s' '%s' '%s'", request.URL.Path, request.RemoteAddr, request.UserAgent())
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
		logwrapper.LogInfof("Error parsing json response")
		return processResult(responseWriter, 3, error.Error())
	}

	return resultJSON
}

func runExecutable(responseWriter http.ResponseWriter, scriptToRun Script) []byte {

	var exitcode int = 0
	var timeoutOccured bool = false
	command := exec.Command(scriptToRun.Path, scriptToRun.Args...)

	if scriptToRun.StdIn != "" {
		command.Stdin = strings.NewReader(scriptToRun.StdIn)
	}

	var timeout = configuration.Settings.Server.DefaultScriptTimeout.Duration
	if scriptToRun.Timeout != "" {
		durationValue, parseError := time.ParseDuration(scriptToRun.Timeout)
		if parseError != nil {
			logwrapper.LogWarningf("Invalid timeout supplied: '%s'", scriptToRun.Timeout)
			return processResult(responseWriter, 3, fmt.Sprintf("Invalid timeout supplied: '%s'", scriptToRun.Timeout))
		}

		timeout = durationValue
	}

	runningProcesses.mutex.Lock()
	runningProcesses.collection[command] = true
	runningProcesses.mutex.Unlock()

	processKiller := time.AfterFunc(timeout, func() {
		command.Process.Kill()
		logwrapper.LogWarningf("Request Timed Out. Timeout: '%s'; Command: '%s'; Arguments: %#v", timeout, scriptToRun.Path, scriptToRun.Args)
		timeoutOccured = true
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

	var outputString = string(output)
	if timeoutOccured {
		exitcode = 3
		outputString = "The script timed out"
	}
	return processResult(responseWriter, exitcode, outputString)
}

func verifySignature(stdin string, signature string) bool {

	stdinAsArray := []byte(stdin)
	signatureStruct, signatureError := minisign.DecodeSignature(signature)

	if signatureError != nil {
		logwrapper.LogInfof("Signature Decoding error: %v", signatureError)
	}

	isValid, error := configuration.Settings.Security.MiniSign.PublicKey.Verify(stdinAsArray, signatureStruct)

	if error != nil {
		logwrapper.LogInfof("Signature Verification Error: %v", error)
	}

	return isValid
}

func verifyRemoteHost(remoteAddr string) bool {

	host, _, parseError := net.SplitHostPort(remoteAddr)

	if parseError != nil {
		logwrapper.LogInfof("Error parsing remote host")
		return false
	}

	remoteIp := net.ParseIP(host)
	allowedAddresses := configuration.Settings.Security.AllowedAddresses.CIDR

	for x := 0; x < len(allowedAddresses); x++ {
		if allowedAddresses[x].Contains(remoteIp) {
			return true
		}
	}
	return false
}

func verifyPathArguments(path string, args []string) bool {
	isValid := false
	allowedSets := configuration.Settings.Security.ApprovedExecutableArguments[path]
	for _, arguments := range allowedSets {
		if cmp.Equal(arguments, args) {
			isValid = true
		}
	}
	return isValid
}
