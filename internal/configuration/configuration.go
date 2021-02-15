package configuration

import (
	"path/filepath"
	"time"

	ini "github.com/vaughan0/go-ini"
)

// SettingsValues is the struct to contain all values
type SettingsValues struct {
	ConfigurationDirectory string
	CertificatePath        string
	PrivateKeyPath         string
	Username               string
	Password               string
	BindAddress            string
	RequestTimeout         time.Duration
}

// Settings is the loaded/updated settings from the configuration file
var Settings = SettingsValues{}

// Initialise loads
func Initialise(configurationDirectory string) {

	Settings.ConfigurationDirectory = configurationDirectory
	Settings.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFile = filepath.FromSlash(configurationDirectory + "/configuration.ini")

	iniFile, loadError := ini.LoadFile(configurationFile)

	if loadError != nil {
		panic(loadError)
	}

	Settings.Username, _ = iniFile.Get("Authentication", "Username")
	Settings.Password, _ = iniFile.Get("Authentication", "Password")

	Settings.BindAddress, _ = iniFile.Get("Server", "BindAddress") // "0.0.0.0:9000",

	stringValue, _ := iniFile.Get("Server", "RequestTimeout")
	Settings.RequestTimeout, _ = time.ParseDuration(stringValue) // "30s"
}
