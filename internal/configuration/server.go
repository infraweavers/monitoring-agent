package configuration

import (
	"encoding/json"
)

// Server is a struct for representing server configuration, unmarshalled from the configuration.json file
type Server struct {
	HTTPRequestTimeout   Duration `json:"HTTPRequestTimeout" mandatory:"true"`
	DefaultScriptTimeout Duration `json:"DefaultScriptTimeout" mandatory:"true"`
	BindAddress          string   `json:"BindAddress" mandatory:"true"`
	LoadPprof            NullBool `json:"LoadPprof" mandatory:"true"`
}

// UnmarshalJSON is a method to implement unmarshalling of the Server type
func (server *Server) UnmarshalJSON(b []byte) error {
	type JSONTmp Server
	var jsonTmp JSONTmp

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
