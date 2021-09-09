package configuration

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/jedisct1/go-minisign"
)

// JSONconfig is a struct for unmarshalling the configuration.json file
type JSONconfig struct {
	Authentication JSONconfigAuthentication `json:"Authentication"`
	Logging        JSONconfigLogging        `json:"Logging"`
	Server         JSONconfigServer         `json:"Server"`
	Security       JSONconfigSecurity       `json:"Security"`
	Paths          JSONconfigPaths          `json:"Paths"`
}

// JSONconfigAuthentication is a struct for unmarshalling the configuration.json file
type JSONconfigAuthentication struct {
	Username string `json:"Username" mandatory:"true"`
	Password string `json:"Password" mandatory:"true"`
}

// JSONconfigLogging is a struct for unmarshalling the configuration.json file, server section
type JSONconfigLogging struct {
	LogFilePath                     string
	LogLevel                        string
	LogArchiveFilesToRetain         int
	LogRotationThresholdInMegaBytes int
	LogHTTPRequests                 bool
	LogHTTPResponses                bool
}

// JSONconfigServer is a struct for unmarshalling the configuration.json file, server section
type JSONconfigServer struct {
	HTTPRequestTimeout           string
	HTTPRequestTimeoutDuration   time.Duration
	DefaultScriptTimeout         string
	DefaultScriptTimeoutDuration time.Duration
	BindAddress                  string
	LoadPprof                    bool
}

// JSONconfigSecurity is a struct for unmarshalling the configuration.json file
type JSONconfigSecurity struct {
	DisableHTTPs                bool
	SignedStdInOnly             bool
	PublicKey                   string
	PublicKeyObject             minisign.PublicKey
	AllowedAddresses            []string
	WhitelistNetworks           []*net.IPNet
	UseClientCertificates       bool
	ClientCertificateCAFile     string
	ClientCertificateCAFilePath string
	ApprovedPathArgumentsOnly   bool
	ApprovedPathArguments       map[string][][]string
}

// JSONconfigPaths is a struct for unmarshalling the configuration.json file
type JSONconfigPaths struct {
	ConfigurationDirectory string
	CertificatePath        string
	PrivateKeyPath         string
}

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	err = validateStruct(v)
	if err != nil {
		return err
	}

	return nil
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

		field := value.Field(i)

		if field.Kind() == reflect.Struct {
			err := validateStruct(field.Interface())
			if err != nil {
				return err
			}
		} else {

			fieldMandatory, _ := strconv.ParseBool(value.Type().Field(i).Tag.Get("mandatory"))

			if fieldMandatory && value.Field(i).IsZero() {
				return fmt.Errorf("%s not set when tagged with 'mandatory:\"true\"'", value.Type().Field(i).Name)
			}
		}
	}
	return nil
}
