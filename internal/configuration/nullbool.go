package configuration

import "encoding/json"

type NullBool struct {
	IsTrue   bool
	HasValue bool
}

func (nullBool *NullBool) UnmarshalJSON(b []byte) error {
	var unmarshalledJson bool

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	nullBool.IsTrue = unmarshalledJson
	nullBool.HasValue = true

	return nil
}
