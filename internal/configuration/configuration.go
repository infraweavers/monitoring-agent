package configuration

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jedisct1/go-minisign"
)

// JSONconfig is a struct for unmarshalling the configuration.json file
type JSONconfig struct {
	Authentication JSONconfigAuthentication `json:"Authentication"`
	Server         JSONconfigServer         `json:"Server"`
	Security       JSONconfigSecurity       `json:"Security"`
	Paths          JSONconfigPaths          `json:"Paths"`
}

// JSONconfigAuthentication is a struct for unmarshalling the configuration.json file
type JSONconfigAuthentication struct {
	Username string
	Password string
}

// JSONconfigServer is a struct for unmarshalling the configuration.json file, server section
type JSONconfigServer struct {
	HTTPRequestTimeout              string
	HTTPRequestTimeoutDuration      time.Duration
	DefaultScriptTimetout           string
	DefaultScriptTimetoutDuration   time.Duration
	LogFilePath                     string
	LogLevel                        string
	BindAddress                     string
	LogArchiveFilesToRetain         int
	LogRotationThresholdInMegaBytes int
	LogHTTPRequests                 bool
	LogHTTPResponses                bool
	LoadPprof                       bool
	DisableHTTPs                    bool
}

// JSONconfigSecurity is a struct for unmarshalling the configuration.json file
type JSONconfigSecurity struct {
	SignedStdInOnly           bool
	PublicKey                 minisign.PublicKey
	AllowedAddresses          []string
	WhitelistNetworks         []*net.IPNet
	UseClientCertificates     bool
	ClientCertificateCAFile   string
	ApprovedPathArgumentsOnly bool
	ApprovedPathArguments     map[string]map[string]bool
}

// JSONconfigPaths is a struct for unmarshalling the configuration.json file
type JSONconfigPaths struct {
	ConfigurationDirectory string
	CertificatePath        string
	PrivateKeyPath         string
}

// Settings is the loaded/updated settings from the configuration file
var Settings = JSONconfig{}

// Initialise loads the settings from the configurationfile
func Initialise(configurationDirectory string) {

	Settings.Paths.ConfigurationDirectory = configurationDirectory
	Settings.Paths.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.Paths.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFileJSON = filepath.FromSlash(configurationDirectory + "/configuration.json")
	var configurationJSON JSONconfig
	jsonFile, jsonErr := ioutil.ReadFile(configurationFileJSON)
	json.Unmarshal(jsonFile, &configurationJSON)
	if jsonErr != nil {
		panic(jsonErr)
	}

	durationValue, parseError := time.ParseDuration(configurationJSON.Server.HTTPRequestTimeout)
	if parseError != nil {
		panic(parseError)
	}
	Settings.Server.HTTPRequestTimeoutDuration = durationValue

	durationValue, parseError = time.ParseDuration(configurationJSON.Server.DefaultScriptTimetout)
	if parseError != nil {
		panic(parseError)
	}
	Settings.Server.DefaultScriptTimetoutDuration = durationValue

	for x := 0; x < len(Settings.Security.AllowedAddresses); x++ {
		_, network, error := net.ParseCIDR(Settings.Security.AllowedAddresses[x])
		if error != nil {
			panic(error)
		}
		Settings.Security.WhitelistNetworks[x] = network
	}

	configurationJSON.Security.ClientCertificateCAFile = fixRelativePath(configurationDirectory, configurationJSON.Security.ClientCertificateCAFile)
}

func fixRelativePath(configurationDirectory string, filePath string) string {
	if filePath == path.Base(filePath) {
		return filepath.FromSlash(configurationDirectory + "/" + filePath)
	}
	return filePath
}

// TestingInitialise only for use in integration tests
func TestingInitialise() {

	configurationDirectoryTemp, _ := os.Getwd()
	configurationDirectory := filepath.FromSlash(configurationDirectoryTemp + "/../../")

	Settings.Paths.ConfigurationDirectory = configurationDirectory
	Settings.Paths.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.Paths.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFileJSON = filepath.FromSlash(configurationDirectory + "/configuration.json")
	jsonFile, jsonErr := ioutil.ReadFile(configurationFileJSON)

	json.Unmarshal(jsonFile, &Settings)

	if jsonErr != nil {
		panic(jsonErr)
	}

	Settings.Server.BindAddress = "127.0.0.1:9000"
	Settings.Server.HTTPRequestTimeoutDuration = time.Second * 11
	Settings.Server.DefaultScriptTimetoutDuration = time.Second * 10

	Settings.Authentication.Username = "test"
	Settings.Authentication.Password = "secret"

	Settings.Security.PublicKey, _ = minisign.NewPublicKey("RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ")

	Settings.Security.WhitelistNetworks = []*net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)},
	}
	Settings.Security.ApprovedPathArgumentsOnly = true
	Settings.Security.ApprovedPathArguments = map[string]map[string]bool{
		`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {
			`-command`: true,
			`-`:        true,
		},
		`sh`: {
			`-c`: true,
			`-s`: true,
		},
	}
}
