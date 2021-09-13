package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Paths is a struct for storing key paths. The ConfigurationDirectory itself is passed in as a flag
// to the main executable or defaulted to the the directory in which the main executable is stored.
// The remaining paths are derived from configuration path.
type Paths struct {
	ConfigurationDirectory string `json:"ConfigurationDirectory" mandatory:"true"`
	CertificatePath        string `json:"CertificatePath" mandatory:"true"`
	PrivateKeyPath         string `json:"PrivateKeyPath" mandatory:"true"`
}

func (paths *Paths) UnmarshalJSON(b []byte) error {
	type JsonTmp Paths
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*paths = Paths(jsonTmp)
	value := reflect.ValueOf(paths)

	for i := 0; i < value.NumField(); i++ {
		isMandatory, _ := strconv.ParseBool(value.Type().Field(i).Tag.Get("mandatory"))
		isZero := value.Field(i).IsZero()

		if isMandatory && isZero {
			return fmt.Errorf("%s not set when tagged with 'mandatory:\"true\"'", value.Type().Field(i).Name)
		}
	}

	return nil
}
