package web

import (
	"encoding/json"
	"fmt"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
	"time"
)

// APIV1RunexecutableGetHandler creates a http response for the API /runexecutable http get requests
func APIV1RunexecutableGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var desc = endpointDescription{
		Endpoint:        "runexecutable",
		Description:     "executes a script as specified in a http request and returns an http response with the result",
		MandatoryFields: "path,args[]",
		OptionalFields:  "timeout",
		ExampleRequest:  `{ "path": "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe", "args":[ "-command", "write-host 'Hello, World'" ] }`,
		ExampleResponse: `{"exitcode":0,"output":"Hello, World\n"}`,
	}
	descJSON, _ := json.Marshal(desc)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(descJSON))
}

// APIV1RunexecutablePostHandler executes a script as specified in a http request and returns an http response with the result
func APIV1RunexecutablePostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := jsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if script.Timeout != "" {
		_, parseError := time.ParseDuration(script.Timeout)
		if parseError != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Invalid timeout supplied: '%s'", http.StatusBadRequest, script.Timeout)))
			logwrapper.LogWarningf("Request received to /runexecutable with invalid timeout supplied: '%s' '%s' '%s'", script.Timeout, request.RemoteAddr, request.UserAgent())
			return
		}
	}

	if configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue {
		if !verifyPathArguments(script.Path, script.Args) {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Unapproved Path/Args", http.StatusBadRequest)))
			logwrapper.LogWarningf("Attempted to use unapproved path and argument(s) combo.", request.RemoteAddr, request.UserAgent())
			return
		}
	}

	responseWriter.Write(runExecutable(responseWriter, script))
}
