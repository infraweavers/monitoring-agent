package configuration

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a struct used for representing time.Duration, unmarshalled from the configuration file
type Duration struct {
	time.Duration
}

// UnmarshalJSON is a method to implement unmarshalling of the Duration type
func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON interface{}

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	switch value := unmarshalledJSON.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration")
	}

	return nil
}
