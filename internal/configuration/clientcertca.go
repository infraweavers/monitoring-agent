package configuration

import (
	"encoding/json"
)

// ClientCertCA is a struct representing Client Certificate Authority path, unmarshalled from the configuration.json file
type ClientCertCA struct {
	Path string
}

// UnmarshalJSON is a method to implement unmarshalling of the ClientCertCA type
func (clientCertCA *ClientCertCA) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON string

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	clientCertCA.Path = fixRelativePath(ConfigurationDirectory, unmarshalledJSON)
	return nil
}
