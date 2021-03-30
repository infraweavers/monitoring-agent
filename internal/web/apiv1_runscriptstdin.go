package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIV1RunscriptstdinGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptstdinGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var desc = endpointDescription{
		Endpoint:        "runscriptstdin",
		Description:     "executes a script included with the post request by passed into a specified command via stdin and returns a http response with the result",
		MandatoryFields: "path,args[],stdin",
		OptionalFields:  "",
		ExampleRequest:  `{ "path": "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe", "args":[ "-command", "-" ], "stdin": "Write-Host 'Hello, World' }`,
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

	if script.StdIn == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint requires stdin", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
