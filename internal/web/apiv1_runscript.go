package web

import (
	"encoding/json"
	"fmt"
	"monitoringagent/internal/configuration"
	"monitoringagent/internal/logwrapper"
	"net/http"
	"strings"
)

// APIV1RunscriptGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var desc = endpointDescription{
		Endpoint:        "runscript",
		Description:     "executes a script as specified in a http request and updates the http response with the result",
		MandatoryFields: "path,args[]",
		OptionalFields:  "",
		ExampleRequest:  `{ "path": "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe", "args":[ "-command", "write-host 'Hello, World'" ] }`,
		ExampleResponse: `{"exitcode":0,"output":"Hello, World\n"}`,
	}
	descJSON, _ := json.Marshal(desc)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(descJSON))
}

// APIV1RunscriptPostHandler executes a script as specified in a http request and updates the http response with the result
func APIV1RunscriptPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := jsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if configuration.Settings.SignedScriptsOnly {
		if script.Signature == "" {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Only signed stdin can be executed", http.StatusBadRequest)))
			return
		}
		if !verifySignature(strings.Join(script.Args[:], " "), script.Signature) {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - Signature not valid", http.StatusBadRequest)))
			logwrapper.Log.Errorf("Attempt to execute script with invalid signature: '%s' '%s' ", request.RemoteAddr, request.UserAgent())
			return
		}
	}

	if script.StdIn != "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint does not use stdin", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
