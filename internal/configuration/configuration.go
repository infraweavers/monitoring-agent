package configuration

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/jedisct1/go-minisign"
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
	LogFilePath            string
	LogLevel               string
	RequestTimeout         time.Duration
	LoadPprof              bool
	SignedStdinOnly        bool
	PublicKey              minisign.PublicKey
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

	Settings.LogFilePath = getIniValueOrPanic(iniFile, "Server", "LogFilePath")
	Settings.LogLevel = getIniValueOrPanic(iniFile, "Server", "LogLevel")

	Settings.BindAddress = getIniValueOrPanic(iniFile, "Server", "BindAddress")

	stringValue := getIniValueOrPanic(iniFile, "Server", "RequestTimeout")
	durationValue, parseError := time.ParseDuration(stringValue)

	if parseError != nil {
		panic(parseError)
	}

	Settings.RequestTimeout = durationValue

	Settings.LoadPprof = getIniBoolOrPanic(iniFile, "Server", "LoadPprof")
	Settings.SignedStdinOnly = getIniBoolOrPanic(iniFile, "Server", "SignedStdinOnly")

	publicKeyString := getIniValueOrPanic(iniFile, "Server", "PublicKey")

	publicKey, publicKeyError := minisign.NewPublicKey(publicKeyString)

	if publicKeyError != nil {
		panic(publicKeyError)
	}
	Settings.PublicKey = publicKey
}

func getIniValueOrPanic(input ini.File, group string, key string) string {
	value, wasFound := input.Get(group, key)
	if wasFound == false {
		panic("[" + group + "] " + key + " was not configured")
	}
	return value
}

func getIniBoolOrPanic(input ini.File, group string, key string) bool {
	stringValue := getIniValueOrPanic(input, group, key)
	boolValue, parseError := strconv.ParseBool(stringValue)

	if parseError != nil {
		panic(parseError)
	}

	return boolValue
}

// TestingInitialise only for use in integration tests
func TestingInitialise() {
	Settings.BindAddress = "127.0.0.1:9000"

	Settings.CertificatePath = "NOTUSED"
	Settings.PrivateKeyPath = "NOTUSED"

	Settings.RequestTimeout = time.Second * 10
	Settings.Username = "test"
	Settings.Password = "secret"
}
