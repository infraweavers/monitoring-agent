package web

import (
	"encoding/json"
	"fmt"
	"monitoringagent/internal/logwrapper"
	"net/http"
)

// APIV1OSSpecificGetHandler creates a http response for the API /runexecutable http get requests
func APIV1OSSpecificGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var desc = endpointDescription{
		Endpoint:        "os_specific",
		Description:     "Handles things that only available on specific Operating Systems (likely Windows)",
		MandatoryFields: "N/A",
		OptionalFields:  "N/A",
		ExampleRequest:  `{ }`,
		ExampleResponse: `{ }`,
	}
	descJSON, _ := json.Marshal(desc)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(descJSON))
}

// APIV1OSSpecificPostHandler creates a http response for the API /info path
func APIV1OSSpecificPostHandler(responseWriter http.ResponseWriter, request *http.Request) {

	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	var requestDecoded OSSpecificRequest
	error := dec.Decode(&requestDecoded)
	if error != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request", http.StatusBadRequest)))
		logwrapper.LogWarningf("Failed JSON decode: '%s' '%s' '%s'", request.URL.Path, request.RemoteAddr, request.UserAgent())
		return
	}
	getValue, err := getResult(requestDecoded)
	if err != nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		logwrapper.LogWarningf("Request received to /os_specific failed : '%s' '%s' '%s'", err, request.RemoteAddr, request.UserAgent())
		return
	}

	resultJSON, _ := json.Marshal(getValue)
	responseWriter.Write(resultJSON)
}
