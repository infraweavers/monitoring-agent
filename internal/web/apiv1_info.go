package web

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// InfoResult internal information about the go runtine
type InfoResult struct {
	Memory runtime.MemStats
}

// APIV1InfoGetHandler creates a http response for the API /info path
func APIV1InfoGetHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reponseWriter.WriteHeader(http.StatusOK)

	response := InfoResult{}
	runtime.ReadMemStats(&response.Memory)

	resultJSON, _ := json.Marshal(response)

	reponseWriter.Write(resultJSON)
}
