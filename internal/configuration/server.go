package configuration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Server is a struct for unmarshalling the configuration.json file, server section
type Server struct {
	HTTPRequestTimeout   Duration `json:"HTTPRequestTimeout" mandatory:"true"`
	DefaultScriptTimeout Duration `json:"DefaultScriptTimeout" mandatory:"true"`
	BindAddress          string   `json:"BindAddress" mandatory:"true"`
	LoadPprof            NullBool `json:"LoadPprof" mandatory:"true"`
}

func (server *Server) UnmarshalJSON(b []byte) error {
	type JsonTmp Server
	var jsonTmp JsonTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*server = Server(jsonTmp)

	value := reflect.ValueOf(server)
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value)
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
