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

// Initialise loads the settings from the configurationfile
func Initialise(configurationDirectory string) {

	Settings.ConfigurationDirectory = configurationDirectory
	Settings.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFile = filepath.FromSlash(configurationDirectory + "/configuration.ini")

	iniFile, loadError := ini.LoadFile(configurationFile)

	if loadError != nil {
		panic(loadError)
	}

	Settings.Username = getIniValueOrPanic(iniFile, "Authentication", "Username")
	Settings.Password = getIniValueOrPanic(iniFile, "Authentication", "Password")

	Settings.BindAddress = getIniValueOrPanic(iniFile, "Server", "BindAddress")

	stringValue := getIniValueOrPanic(iniFile, "Server", "RequestTimeout")
	durationValue, parseError := time.ParseDuration(stringValue)

	if parseError != nil {
		panic(parseError)
	}

	Settings.RequestTimeout = durationValue
}

func getIniValueOrPanic(input ini.File, group string, key string) string {
	value, wasFound := input.Get(group, key)
	if wasFound == false {
		panic("[" + group + "] " + key + " was not configured")
	}
	return value
}
