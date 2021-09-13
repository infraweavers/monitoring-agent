package configuration

import (
	"encoding/json"
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

	err = validateStruct(server)
	if err != nil {
		return err
	}

	return nil
}
