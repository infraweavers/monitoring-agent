package configuration

import (
	"encoding/json"
	"path/filepath"
)

// Paths is a struct for storing key paths. The ConfigurationDirectory itself is passed in as a flag
// to the main executable or defaulted to the the directory in which the main executable is stored.
// The remaining paths are derived from configuration path.
type Paths struct {
	ConfigurationDirectory string `json:"-"`
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

func InitialisePaths(configurationDirectory string) Paths {
	return Paths{
		ConfigurationDirectory: configurationDirectory,
		CertificatePath:        filepath.FromSlash(configurationDirectory + "/server.crt"),
		PrivateKeyPath:         filepath.FromSlash(configurationDirectory + "/server.key"),
	}
}

func (paths *Paths) mmmmm(p Paths) {
	paths.ConfigurationDirectory = p.ConfigurationDirectory

	if Settings.Paths.CertificatePath == "" {
		Settings.Paths.CertificatePath = paths.CertificatePath
	}

	if Settings.Paths.PrivateKeyPath == "" {
		Settings.Paths.PrivateKeyPath = paths.PrivateKeyPath
	}
}
