package configuration

import (
	"encoding/json"
)

// Authentication is a struct for unmarshalling the configuration.json file
type Authentication struct {
	Username string `json:"Username" mandatory:"true"`
	Password string `json:"Password" mandatory:"true"`
}

func (authentication *Authentication) UnmarshalJSON(b []byte) error {
	type JsonTmp Authentication
	var jsonTmp JsonTmp

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
