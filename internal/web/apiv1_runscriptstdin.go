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
		Description:     "executes a script included with the post request by passed into a specified command via stdin and returns a http response with the result",
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

// APIV1RunscriptstdinPostHandler creates executes a script by piping it to the standard input (stdin) of the specified command
func APIV1RunscriptstdinPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := jsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if configuration.Settings.SignedStdInOnly {
		if script.StdInSignature == "" {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Only signed stdin can be executed", http.StatusBadRequest)))
			return
		}
		if !verifySignature(script.StdIn, script.StdInSignature) {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Signature not valid", http.StatusBadRequest)))
			logwrapper.Log.Errorf("Attempt to execute script with invalid signature: '%s' '%s' ", request.RemoteAddr, request.UserAgent())
			return
		}
	}

	if script.StdIn == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint requires stdin", http.StatusBadRequest)))
		logwrapper.Log.Errorf("Attempt to execute script as stdin without stdin in the request")
		return
	}

	if script.Timeout != "" {
		_, parseError := time.ParseDuration(script.Timeout)
		if parseError != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Invalid timeout supplied: '%s'", http.StatusBadRequest, script.Timeout)))
			return
		}
	}

	responseWriter.Write(runScript(responseWriter, script))
}
