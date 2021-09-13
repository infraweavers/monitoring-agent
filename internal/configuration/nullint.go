package configuration

import "encoding/json"

type NullInt struct {
	Value    int
	HasValue bool
}

func (nullInt *NullInt) UnmarshalJSON(b []byte) error {
	var unmarshalledJson int

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	nullInt.Value = unmarshalledJson
	nullInt.HasValue = true

	return nil
}
