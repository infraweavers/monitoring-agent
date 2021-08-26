package configuration

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jedisct1/go-minisign"
)

// SettingsValues is the struct to contain all values
type SettingsValues struct {
	ConfigurationDirectory          string
	CertificatePath                 string
	PrivateKeyPath                  string
	Username                        string
	Password                        string
	BindAddress                     string
	LogFilePath                     string
	LogLevel                        string
	LogArchiveFilesToRetain         int
	LogRotationThresholdInMegaBytes int
	LogHTTPRequests                 bool
	LogHTTPResponses                bool
	HTTPRequestTimeout              time.Duration
	DefaultScriptTimeout            time.Duration
	LoadPprof                       bool
	DisableHTTPs                    bool
	SignedStdInOnly                 bool
	PublicKey                       minisign.PublicKey
	AllowedAddresses                []*net.IPNet
	UseClientCertificates           bool
	ClientCertificateCAFile         string
	ApprovedPathArgumentsOnly       bool
	ApprovedPathArguments           map[string]map[string]bool
}

// struct for unmarshalling the configuration.json file
type JSONconfig struct {
	Authentication JSONconfigAuthentication `json:"Authentication"`
	Server         JSONconfigServer         `json:"Server"`
	Security       JSONconfigSecurity       `json:"Security"`
}

// struct for unmarshalling the configuration.json file
type JSONconfigAuthentication struct {
	Username string
	Password string
}

// struct for unmarshalling the configuration.json file
type JSONconfigServer struct {
	HTTPRequestTimeout              string
	DefaultScriptTimetout           string
	logFilePath                     string
	LogLevel                        string
	BindAddress                     string
	LogArchiveFilesToRetain         int
	LogRotationThresholdInMegaBytes int
	LogHTTPRequests                 bool
	LogHTTPResponses                bool
	LoadPprof                       bool
	DisabledHTTPs                   bool
}

// struct for unmarshalling the configuration.json file
type JSONconfigSecurity struct {
	SignedStdInOnly           bool
	PublicKey                 string
	AllowedAddresses          string
	UseClientCertificates     bool
	ClientCertificateCAFile   string
	ApprovedPathArgumentsOnly bool
	ApprovedPathArguments     map[string]map[string]bool
}

// Settings is the loaded/updated settings from the configuration file
var Settings = SettingsValues{}

// Initialise loads the settings from the configurationfile
func Initialise(configurationDirectory string) {

	Settings.ConfigurationDirectory = configurationDirectory

	Settings.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFileJSON = filepath.FromSlash(configurationDirectory + "/configuration.json")

	var configurationJSON JSONconfig

	jsonFile, jsonErr := ioutil.ReadFile(configurationFileJSON)

	json.Unmarshal(jsonFile, &configurationJSON)

	if jsonErr != nil {
		panic(jsonErr)
	}

	Settings.Username = configurationJSON.Authentication.Username
	Settings.Password = configurationJSON.Authentication.Password

	Settings.LogFilePath = configurationJSON.Server.logFilePath
	Settings.LogLevel = configurationJSON.Server.LogLevel
	Settings.LogArchiveFilesToRetain = configurationJSON.Server.LogArchiveFilesToRetain
	Settings.LogRotationThresholdInMegaBytes = configurationJSON.Server.LogRotationThresholdInMegaBytes
	Settings.LogHTTPRequests = configurationJSON.Server.LogHTTPRequests
	Settings.LogHTTPResponses = configurationJSON.Server.LogHTTPResponses
	Settings.BindAddress = configurationJSON.Server.BindAddress

	durationValue, parseError := time.ParseDuration(configurationJSON.Server.HTTPRequestTimeout)
	if parseError != nil {
		panic(parseError)
	}

	Settings.HTTPRequestTimeout = durationValue

	durationValue, parseError = time.ParseDuration(configurationJSON.Server.DefaultScriptTimetout)
	if parseError != nil {
		panic(parseError)
	}
	Settings.DefaultScriptTimeout = durationValue

	Settings.LoadPprof = configurationJSON.Server.LoadPprof
	Settings.DisableHTTPs = configurationJSON.Server.DisabledHTTPs

	Settings.SignedStdInOnly = configurationJSON.Security.SignedStdInOnly

	hostArrays := strings.Split(configurationJSON.Security.AllowedAddresses, ",")
	whitelistNetworks := make([]*net.IPNet, len(hostArrays))
	for x := 0; x < len(hostArrays); x++ {
		_, network, error := net.ParseCIDR(hostArrays[x])
		if error != nil {
			panic(error)
		}
		whitelistNetworks[x] = network
	}
	Settings.AllowedAddresses = whitelistNetworks

	publicKeyString := configurationJSON.Security.PublicKey
	publicKey, publicKeyError := minisign.NewPublicKey(publicKeyString)

	if publicKeyError != nil {
		panic(publicKeyError)
	}
	Settings.PublicKey = publicKey

	Settings.UseClientCertificates = configurationJSON.Security.UseClientCertificates
	Settings.ClientCertificateCAFile = fixRelativePath(configurationDirectory, configurationJSON.Security.ClientCertificateCAFile)
}

func fixRelativePath(configurationDirectory string, filePath string) string {
	if filePath == path.Base(filePath) {
		return filepath.FromSlash(configurationDirectory + "/" + filePath)
	}
	return filePath
}

// TestingInitialise only for use in integration tests
func TestingInitialise() {

	// TESTING CONFIG FILES SECTION
	//configurationDirectory := `D:\code\monitoring-agent`
	configurationDirectoryTemp, _ := os.Getwd()
	configurationDirectory := filepath.FromSlash(configurationDirectoryTemp + "/../../")

	Settings.ConfigurationDirectory = filepath.FromSlash(configurationDirectory + "/../../")

	Settings.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFileJSON = filepath.FromSlash(configurationDirectory + "/configuration.json")

	// JSON
	var configurationJSON JSONconfig

	jsonFile, jsonErr := ioutil.ReadFile(configurationFileJSON)

	json.Unmarshal(jsonFile, &configurationJSON)

	if jsonErr != nil {
		panic(jsonErr)
	}
	// END TESTING CONFIG FILES SECTION

	Settings.BindAddress = "127.0.0.1:9000"

	Settings.CertificatePath = "NOTUSED"
	Settings.PrivateKeyPath = "NOTUSED"

	Settings.HTTPRequestTimeout = time.Second * 11
	Settings.DefaultScriptTimeout = time.Second * 10
	Settings.Username = "test"
	Settings.Password = "secret"

	Settings.PublicKey, _ = minisign.NewPublicKey("RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ")

	Settings.AllowedAddresses = []*net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)},
	}
	Settings.ApprovedPathArgumentsOnly = true
	Settings.ApprovedPathArguments = map[string]map[string]bool{
		`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {
			`-command`: true,
			`-`:        true,
		},
		`sh`: {
			`-c`: true,
			`-s`: true,
		},
	}
	Settings.ApprovedPathArguments = configurationJSON.Security.ApprovedPathArguments
}
