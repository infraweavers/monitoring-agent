package apiv1

import (
	"net/http"
	"encoding/json"
	"os/exec"
	"log"
	"internal/basicauth"
)

type Script struct {
	Path	string		`json:"path"`
	Args	[]string	`json:"args"`
}

type Result struct {
	Exitcode	int		`json:"exitcode"`
	Output		string	`json:"output"`
}

func runsScript(s Script) Result {

	var exitcode int = 0
	log.Printf("Executing: %+v", s)
	
	cmd := exec.Command( s.Path, s.Args... )
	output, err := cmd.CombinedOutput()
    if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		}
        log.Println( "Error:", err )
    } 
	log.Printf("Result: %s", output)
	
	result := Result{
		Exitcode: exitcode,
		Output: string(output),
	}
	
	return result
}

func RunscriptGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	basicauth.IsAuthorised(w, r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"endpoints": ["executes a script sent via HTTP POST request"]}`))
}

func RunscriptPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	basicauth.IsAuthorised(w, r)
	
	dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
	var script Script
	err := dec.Decode(&script)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	
	result := runsScript(script)
	json.NewEncoder(w).Encode(result)
}