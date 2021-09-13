package configuration

import (
	"encoding/json"

	"github.com/jedisct1/go-minisign"
)

type MiniSign struct {
	minisign.PublicKey
}

func (ms *MiniSign) UnmarshalJSON(b []byte) error {
	var unmarshalledJson string

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	ms.PublicKey, err = minisign.NewPublicKey(unmarshalledJson)
	if err != nil {
		return err
	}

	return nil
}
