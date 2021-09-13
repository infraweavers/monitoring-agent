package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Authentication is a struct for unmarshalling the configuration.json file
type Authentication struct {
	Username string `json:"Username" mandatory:"true"`
	Password string `json:"Password" mandatory:"true"`
}

func (authentication *Authentication) UnmarshalJSON(b []byte) error {
	type JsonTmp Authentication
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*authentication = Authentication(jsonTmp)

	value := reflect.ValueOf(authentication)
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
