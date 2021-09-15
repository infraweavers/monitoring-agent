package configuration

import (
	"encoding/json"
)

// Authentication is a struct representing credentials for using monitoring-agent, which is unmarshalled from the configuration.json file
type Authentication struct {
	Username string `json:"Username" mandatory:"true"`
	Password string `json:"Password" mandatory:"true"`
}

// UnmarshalJSON is a method to implement unmarshalling of the Authentication type
func (authentication *Authentication) UnmarshalJSON(b []byte) error {
	type JSONTmp Authentication
	var jsonTmp JSONTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*authentication = Authentication(jsonTmp)

	err = validateStruct(authentication)
	if err != nil {
		return err
	}

	return nil
}
