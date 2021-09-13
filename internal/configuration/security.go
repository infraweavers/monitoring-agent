package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JSONconfigSecurity is a struct for unmarshalling the configuration.json file
type Security struct {
	DisableHTTPs              NullBool              `json:"DisableHTTPs" mandatory:"true"`
	SignedStdInOnly           NullBool              `json:"SignedStdInOnly" mandatory:"true"`
	MiniSign                  MiniSign              `json:"PublicKey" mandatory:"true"`
	AllowedAddresses          AllowedNetworks       `json:"AllowedAddresses" mandatory:"true"`
	UseClientCertificates     NullBool              `json:"UseClientCertificates" mandatory:"true"`
	ClientCertificateCAFile   ClientCertCA          `json:"ClientCertificateCAFile" mandatory:"true"`
	ApprovedPathArgumentsOnly NullBool              `json:"ApprovedPathArgumentsOnly" mandatory:"true"`
	ApprovedPathArguments     map[string][][]string `json:"ApprovedPathArguments" mandatory:"true"`
}

func (security *Security) UnmarshalJSON(b []byte) error {
	type JsonTmp Security
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*security = Security(jsonTmp)

	value := reflect.ValueOf(security)
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
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
