package apiv1

import (
	"encoding/json"
	"log"
	"mama/internal/basicauth"
	"net/http"
	"os/exec"
)

type script struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type result struct {
	Exitcode int    `json:"exitcode"`
	Output   string `json:"output"`
}

func runsScript(scriptToRun script) result {

	var exitcode int = 0
	log.Printf("Executing: %+v", scriptToRun)

	cmd := exec.Command(scriptToRun.Path, scriptToRun.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		}
		log.Println("Error:", err)
	}
	log.Printf("Result: %s", output)

	scriptResult := result{
		Exitcode: exitcode,
		Output:   string(output),
	}

	return scriptResult
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
	var script script
	err := dec.Decode(&script)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	result := runsScript(script)
	resultJson, err := json.Marshal(result)
	responseWriter.Write(resultJson)
}
