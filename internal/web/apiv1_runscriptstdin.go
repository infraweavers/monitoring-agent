package web

import (
	"encoding/json"
	"fmt"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
	"time"
)

// APIV1RunscriptstdinGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptstdinGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var desc = endpointDescription{
		Endpoint:        "runscriptstdin",
		Description:     "executes a script included with the post request, passing into the specified command via stdin and returns a http response with the result",
		MandatoryFields: "path,args[],stdin",
		OptionalFields:  "stdinsignature, timeout",
		ExampleRequest:  `{ "path": "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe", "args":[ "-command", "-" ], "stdin": "Write-Host 'Hello, World'" }`,
		ExampleResponse: `{"exitcode":0,"output":"Hello, World\n"}`,
	}
	descJSON, _ := json.Marshal(desc)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(descJSON))
}

// APIV1RunscriptstdinPostHandler executes a script by piping it to the standard input (stdin) of the specified command
func APIV1RunscriptstdinPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := jsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if configuration.Settings.Security.ApprovedPathArgumentsOnly {
		if !verifyPathArguments(script.Path, script.Args) {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Unapproved Path/Args", http.StatusBadRequest)))
			logwrapper.LogWarningf("Attempted to use unapproved path and argument(s) combo.", request.RemoteAddr, request.UserAgent())
			return
		}
	}

	if configuration.Settings.Security.SignedStdInOnly {
		if script.StdInSignature == "" {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Only signed stdin can be executed", http.StatusBadRequest)))
			logwrapper.LogWarningf("Attempt to execute script with no signature when configuration.Settings.SignedStdInOnly is %t: '%s' '%s' ", configuration.Settings.Security.SignedStdInOnly, request.RemoteAddr, request.UserAgent())
			return
		}
		if !verifySignature(script.StdIn, script.StdInSignature) {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Signature not valid", http.StatusBadRequest)))
			logwrapper.LogWarningf("Attempt to execute script with invalid signature: '%s' '%s' ", request.RemoteAddr, request.UserAgent())
			return
		}
	}

	if script.StdIn == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint requires stdin", http.StatusBadRequest)))
		logwrapper.LogWarningf("Attempt to execute script as stdin without stdin in the request: '%s' '%s' ", request.RemoteAddr, request.UserAgent())
		return
	}

	if script.Timeout != "" {
		_, parseError := time.ParseDuration(script.Timeout)
		if parseError != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Invalid timeout supplied: '%s'", http.StatusBadRequest, script.Timeout)))
			logwrapper.LogWarningf("Request received to /runscriptstdin with invalid timeout supplied: '%s' '%s' '%s'", script.Timeout, request.RemoteAddr, request.UserAgent())
			return
		}
	}

	responseWriter.Write(runScript(responseWriter, script))
}
