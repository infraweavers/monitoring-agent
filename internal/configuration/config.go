package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JsonConfig is a struct for unmarshalling the configuration.json file
type Config struct {
	Authentication Authentication `json:"Authentication" mandatory:"true"`
	Logging        Logging        `json:"Logging" mandatory:"true"`
	Server         Server         `json:"Server" mandatory:"true"`
	Security       Security       `json:"Security" mandatory:"true"`
	Paths          Paths          `json:"Paths"`
}

func (config *Config) UnmarshalJSON(b []byte) error {
	type JsonTmp Config
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*config = Config(jsonTmp)

	value := reflect.ValueOf(config)
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
