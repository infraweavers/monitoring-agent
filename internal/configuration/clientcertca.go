package configuration

import (
	"encoding/json"
)

type ClientCertCA struct {
	Path string
}

func (clientCertCA *ClientCertCA) UnmarshalJSON(b []byte) error {
	var unmarshalledJson string

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	clientCertCA.Path = fixRelativePath(Settings.Paths.ConfigurationDirectory, unmarshalledJson)
	return nil
}
