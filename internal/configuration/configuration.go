package configuration

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jedisct1/go-minisign"
)

// Settings is the loaded/updated settings from the configuration file
var Settings = JSONconfig{}

// Initialise loads the settings from the configurationfile
func Initialise(configurationDirectory string) {

	Settings.Paths.ConfigurationDirectory = configurationDirectory
	Settings.Paths.CertificatePath = filepath.FromSlash(configurationDirectory + "/server.crt")
	Settings.Paths.PrivateKeyPath = filepath.FromSlash(configurationDirectory + "/server.key")

	var configurationFileJSON = filepath.FromSlash(configurationDirectory + "/configuration.json")
	jsonFile, err := ioutil.ReadFile(configurationFileJSON)
	if err != nil {
		panic(err)
	}

	err = Unmarshal(jsonFile, &Settings)
	if err != nil {
		panic(err)
	}

	for x := 0; x < len(Settings.Security.AllowedAddresses); x++ {
		_, network, error := net.ParseCIDR(Settings.Security.AllowedAddresses[x])
		if error != nil {
			panic(error)
		}
		Settings.Security.WhitelistNetworks = append(Settings.Security.WhitelistNetworks, network)
	}

	publicKeyObject, error := minisign.NewPublicKey(Settings.Security.PublicKey)
	if error != nil {
		panic(error)
	}
	Settings.Security.PublicKeyObject = publicKeyObject
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

	Initialise(configurationDirectory)

	Settings.Server.BindAddress = "127.0.0.1:9000"
	Settings.Server.HTTPRequestTimeout.Duration = time.Second * 11
	Settings.Server.DefaultScriptTimeout.Duration = time.Second * 10

	Settings.Authentication.Username = "test"
	Settings.Authentication.Password = "secret"

	Settings.Security.PublicKeyObject, _ = minisign.NewPublicKey("RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ")

	Settings.Security.WhitelistNetworks = []*net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)},
	}
}
