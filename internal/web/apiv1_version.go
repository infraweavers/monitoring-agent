package web

import (
	"encoding/json"
	"monitoringagent/internal/configuration"
	"net/http"
)

type VersionResult struct {
	Version string
}

// APIV1InfoGetHandler creates a http response for the API /info path
func APIV1VersionHandler(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reponseWriter.WriteHeader(http.StatusOK)

	response := VersionResult{}
	response.Version = configuration.Settings.MonitoringAgentVersion

	resultJSON, _ := json.Marshal(response)
	reponseWriter.Write(resultJSON)
}
