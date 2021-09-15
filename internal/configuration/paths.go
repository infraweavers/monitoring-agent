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

// UnmarshalJSON is a method to implement unmarshalling of the Paths type
func (paths *Paths) UnmarshalJSON(b []byte) error {
	type JSONTmp Paths
	var jsonTmp JSONTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*paths = Paths(jsonTmp)
	err = validateStruct(paths)
	if err != nil {
		return err
	}

	paths.ConfigurationDirectory = ConfigurationDirectory
	return nil
}

// InitialisePaths is a function for setting Paths.CertificatePath and Paths.PrivateKeyPath to their default values based on ConfigurationDirectory
func InitialisePaths(configurationDirectory string) Paths {
	ConfigurationDirectory = configurationDirectory

	return Paths{
		ConfigurationDirectory: ConfigurationDirectory,
		CertificatePath:        filepath.FromSlash(configurationDirectory + "/server.crt"),
		PrivateKeyPath:         filepath.FromSlash(configurationDirectory + "/server.key"),
	}
}

// Reset is a method for resetting Paths.CertificatePath and Paths.PrivateKeyPath to their default values based on ConfigurationDirectory, if
// they remain as empty string ("") after importing the configuration.json file
func (paths *Paths) Reset(p Paths) {
	paths.ConfigurationDirectory = p.ConfigurationDirectory

	if Settings.Paths.CertificatePath == "" {
		Settings.Paths.CertificatePath = p.CertificatePath
	}

	if Settings.Paths.PrivateKeyPath == "" {
		Settings.Paths.PrivateKeyPath = p.PrivateKeyPath
	}
}
