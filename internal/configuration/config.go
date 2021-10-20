package configuration

import (
	"encoding/json"
)

// Config is a struct for unmarshalling the configuration.json file
type Config struct {
	Authentication         Authentication `json:"Authentication" mandatory:"true"`
	Logging                Logging        `json:"Logging" mandatory:"true"`
	Server                 Server         `json:"Server" mandatory:"true"`
	Security               Security       `json:"Security" mandatory:"true"`
	Paths                  Paths          `json:"Paths"`
	MonitoringAgentVersion string
}

// UnmarshalJSON is a method to implement unmarshalling of the Config type
func (config *Config) UnmarshalJSON(b []byte) error {
	type JSONTmp Config
	var jsonTmp JSONTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*config = Config(jsonTmp)

	err = validateStruct(config)
	if err != nil {
		return err
	}

	return nil
}
