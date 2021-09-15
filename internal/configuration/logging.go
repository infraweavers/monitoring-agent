package configuration

import (
	"encoding/json"
)

// Logging is a struct representing the logging configuration, unmarshalled from the configuration.json file
type Logging struct {
	LogFilePath                     string   `json:"LogFilePath" mandatory:"true"`
	LogLevel                        string   `json:"LogLevel" mandatory:"true"`
	LogArchiveFilesToRetain         NullInt  `json:"LogArchiveFilesToRetain" mandatory:"true"`
	LogRotationThresholdInMegaBytes int      `json:"LogRotationThresholdInMegaBytes" mandatory:"true"`
	LogHTTPRequests                 NullBool `json:"LogHTTPRequests" mandatory:"true"`
	LogHTTPResponses                NullBool `json:"LogHTTPResponses" mandatory:"true"`
}

// UnmarshalJSON is a method to implement unmarshalling of the AllowedNetworks type
func (logging *Logging) UnmarshalJSON(b []byte) error {
	type JSONTmp Logging
	var jsonTmp JSONTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*logging = Logging(jsonTmp)
	logging.LogFilePath = fixRelativePath(ConfigurationDirectory, logging.LogFilePath)

	err = validateStruct(logging)
	if err != nil {
		return err
	}

	return nil
}
