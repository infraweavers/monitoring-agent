package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/jedisct1/go-minisign"
)

// Settings is the loaded/updated settings from the configuration file
var Settings = Config{}

// Initialise loads the settings from the configurationfile
func Initialise(configurationDirectory string) {

	var paths = Paths{
		ConfigurationDirectory: configurationDirectory,
		CertificatePath:        filepath.FromSlash(configurationDirectory + "/server.crt"),
		PrivateKeyPath:         filepath.FromSlash(configurationDirectory + "/server.key"),
	}

	var configurationFileJSON = filepath.FromSlash(paths.ConfigurationDirectory + "/configuration.json")
	jsonFile, err := ioutil.ReadFile(configurationFileJSON)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonFile, &Settings)
	if err != nil {
		panic(err)
	}

	Settings.Paths.ConfigurationDirectory = paths.ConfigurationDirectory

	if Settings.Paths.CertificatePath == "" {
		Settings.Paths.CertificatePath = paths.CertificatePath
	}

	if Settings.Paths.PrivateKeyPath == "" {
		Settings.Paths.PrivateKeyPath = paths.PrivateKeyPath
	}

	Settings.Logging.LogFilePath = fixRelativePath(configurationDirectory, Settings.Logging.LogFilePath)
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

	Settings.Security.MiniSign.PublicKey, _ = minisign.NewPublicKey("RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ")

	Settings.Security.AllowedAddresses.CIDR = []*net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)},
	}
}

func fixRelativePath(configurationDirectory string, filePath string) string {
	if filePath == path.Base(filePath) {
		return filepath.FromSlash(configurationDirectory + "/" + filePath)
	}
	return filePath
}

func validateStruct(item interface{}) error {

	value := reflect.ValueOf(item)

	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value)
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("value type was %s rather than struct", reflect.TypeOf(value))
	}

	for i := 0; i < value.NumField(); i++ {
		isMandatory, _ := strconv.ParseBool(value.Type().Field(i).Tag.Get("mandatory"))
		isZero := value.Field(i).IsZero()

		if isMandatory && isZero {
			return fmt.Errorf("%s not set when tagged with 'mandatory:\"true\"'", value.Type().Field(i).Name)
		}
	}
	return nil
}
