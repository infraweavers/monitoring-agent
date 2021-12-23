package web

import (
	"encoding/json"
	"net/http"
)

type CounterResult struct {
	Name  string
	Value string
}

// APIV1OSSpecificHandler creates a http response for the API /info path
func APIV1OSSpecificHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reponseWriter.WriteHeader(http.StatusOK)

	counter := "\\Memory\\Available MBytes"

	getValue := getResult(counter)

	resultJSON, _ := json.Marshal(getValue)
	reponseWriter.Write(resultJSON)
}
