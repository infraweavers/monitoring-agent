package configuration

import (
	"encoding/json"

	"github.com/jedisct1/go-minisign"
)

// MiniSign is a struct representing a minisign.Public key used for verifying signed scripts, unmarshalled from the configuration.json file
type MiniSign struct {
	minisign.PublicKey
}

// UnmarshalJSON is a method to implement unmarshalling of the MiniSign type
func (ms *MiniSign) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON string

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	ms.PublicKey, err = minisign.NewPublicKey(unmarshalledJSON)
	if err != nil {
		return err
	}

	return nil
}
