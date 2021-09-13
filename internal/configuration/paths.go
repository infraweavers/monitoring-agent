package configuration

import (
	"encoding/json"
)

// Paths is a struct for storing key paths. The ConfigurationDirectory itself is passed in as a flag
// to the main executable or defaulted to the the directory in which the main executable is stored.
// The remaining paths are derived from configuration path.
type Paths struct {
	ConfigurationDirectory string `json:"ConfigurationDirectory"`
	CertificatePath        string `json:"CertificatePath" mandatory:"true"`
	PrivateKeyPath         string `json:"PrivateKeyPath" mandatory:"true"`
}

func (paths *Paths) UnmarshalJSON(b []byte) error {
	type JsonTmp Paths
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*paths = Paths(jsonTmp)
	err = validateStruct(paths)
	if err != nil {
		return err
	}

	return nil
}
