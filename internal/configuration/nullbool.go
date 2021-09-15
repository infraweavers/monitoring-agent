package configuration

import "encoding/json"

// NullBool is a struct giving a distinction between a bool that is false because it is unset and a bool that has explicitly been set as false.
// Used for representing booleans unmarshalled from the configuration.json file
type NullBool struct {
	IsTrue   bool
	HasValue bool
}

// UnmarshalJSON is a method to implement unmarshalling of the NullBool type
func (nullBool *NullBool) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON bool

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	nullBool.IsTrue = unmarshalledJSON
	nullBool.HasValue = true

	return nil
}
