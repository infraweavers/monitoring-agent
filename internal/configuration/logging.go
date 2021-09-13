package configuration

import (
	"encoding/json"
)

// JSONconfigLogging is a struct for unmarshalling the configuration.json file, server section
type Logging struct {
	LogFilePath                     string   `json:"LogFilePath" mandatory:"true"`
	LogLevel                        string   `json:"LogLevel" mandatory:"true"`
	LogArchiveFilesToRetain         int      `json:"LogArchiveFilesToRetain" mandatory:"true"`
	LogRotationThresholdInMegaBytes int      `json:"LogRotationThresholdInMegaBytes" mandatory:"true"`
	LogHTTPRequests                 NullBool `json:"LogHTTPRequests" mandatory:"true"`
	LogHTTPResponses                NullBool `json:"LogHTTPResponses" mandatory:"true"`
}

func (logging *Logging) UnmarshalJSON(b []byte) error {
	type JsonTmp Logging
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*logging = Logging(jsonTmp)

	err = validateStruct(logging)
	if err != nil {
		return err
	}

	return nil
}
