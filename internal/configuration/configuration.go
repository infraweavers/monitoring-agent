package configuration

import (
	"net"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jedisct1/go-minisign"
	ini "github.com/vaughan0/go-ini"
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
	ApprovedExecutable              bool
	ApprovedPath                    map[string]bool
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

	Settings.LogFilePath = fixRelativePath(configurationDirectory, getIniValueOrPanic(iniFile, "Server", "LogFilePath"))
	Settings.LogLevel = getIniValueOrPanic(iniFile, "Server", "LogLevel")

	stringValue := getIniValueOrPanic(iniFile, "Server", "LogArchiveFilesToRetain")
	intValue, parseError := strconv.Atoi(stringValue)
	if parseError != nil {
		panic(parseError)
	}
	Settings.LogArchiveFilesToRetain = intValue

	stringValue = getIniValueOrPanic(iniFile, "Server", "LogRotationThresholdInMegaBytes")
	intValue, parseError = strconv.Atoi(stringValue)
	if parseError != nil {
		panic(parseError)
	}
	Settings.LogRotationThresholdInMegaBytes = intValue

	Settings.LogHTTPRequests = getIniBoolOrPanic(iniFile, "Server", "LogHTTPRequests")
	Settings.LogHTTPResponses = getIniBoolOrPanic(iniFile, "Server", "LogHTTPResponses")

	Settings.BindAddress = getIniValueOrPanic(iniFile, "Server", "BindAddress")

	stringValue = getIniValueOrPanic(iniFile, "Server", "HTTPRequestTimeout")
	durationValue, parseError := time.ParseDuration(stringValue)
	if parseError != nil {
		panic(parseError)
	}
	Settings.HTTPRequestTimeout = durationValue

	stringValue = getIniValueOrPanic(iniFile, "Server", "DefaultScriptTimeout")
	durationValue, parseError = time.ParseDuration(stringValue)
	if parseError != nil {
		panic(parseError)
	}
	Settings.DefaultScriptTimeout = durationValue

	Settings.LoadPprof = getIniBoolOrPanic(iniFile, "Server", "LoadPprof")
	Settings.DisableHTTPs = getIniBoolOrPanic(iniFile, "Server", "DisableHTTPs")

	Settings.SignedStdInOnly = getIniBoolOrPanic(iniFile, "Security", "SignedStdInOnly")

	hostArrays := strings.Split(getIniValueOrPanic(iniFile, "Security", "AllowedAddresses"), ",")
	whitelistNetworks := make([]*net.IPNet, len(hostArrays))
	for x := 0; x < len(hostArrays); x++ {
		_, network, error := net.ParseCIDR(hostArrays[x])
		if error != nil {
			panic(error)
		}
		whitelistNetworks[x] = network
	}
	Settings.AllowedAddresses = whitelistNetworks

	publicKeyString := getIniValueOrPanic(iniFile, "Security", "PublicKey")
	publicKey, publicKeyError := minisign.NewPublicKey(publicKeyString)

	if publicKeyError != nil {
		panic(publicKeyError)
	}
	Settings.PublicKey = publicKey

	Settings.UseClientCertificates = getIniBoolOrPanic(iniFile, "Security", "UseClientCertificates")
	Settings.ClientCertificateCAFile = fixRelativePath(configurationDirectory, getIniValueOrPanic(iniFile, "Security", "ClientCertificateCAFile"))
}

func getIniValueOrPanic(input ini.File, group string, key string) string {
	value, wasFound := input.Get(group, key)
	if !wasFound {
		panic("[" + group + "] " + key + " was not configured")
	}
	return value
}

func fixRelativePath(configurationDirectory string, filePath string) string {
	if filePath == path.Base(filePath) {
		return filepath.FromSlash(configurationDirectory + "/" + filePath)
	}
	return filePath
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

	Settings.HTTPRequestTimeout = time.Second * 11
	Settings.DefaultScriptTimeout = time.Second * 10
	Settings.Username = "test"
	Settings.Password = "secret"

	Settings.PublicKey, _ = minisign.NewPublicKey("RWQ3ly9IPenQ6Wgt/VYzMCdGdVJPPoNSyT+rtTddvqBgANTYdboko0zu")
	Settings.PublicKey, _ = minisign.NewPublicKey("RWTVYlcv8rHLCPg9ME+2wyEtwHz1azX54uLnGW5AWzb1R1qaESVNzxGI")
	Settings.AllowedAddresses = []*net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)},
	}
	Settings.ApprovedExecutable = false
	Settings.ApprovedPath = map[string]bool{
		`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: true,
		`sh`: true,
	}
}
