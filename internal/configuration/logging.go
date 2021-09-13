package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

	value := reflect.ValueOf(logging)
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		isMandatory, _ := strconv.ParseBool(value.Type().Field(i).Tag.Get("mandatory"))
		isZero := value.Field(i).IsZero()

		if isMandatory && isZero {
			return fmt.Errorf("%s not set when tagged with 'mandatory:\"true\"'", value.Type().Field(i).Name)
		}
	}

	return nil
}
