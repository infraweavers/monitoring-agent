package configuration

import "encoding/json"

// NullInt is a struct giving a distinction between an int that is 0 because it is unset and an int that has explicitly been set as 0.
// Used for representing integers unmarshalled from the configuration.json file
type NullInt struct {
	Value    int
	HasValue bool
}

// UnmarshalJSON is a method to implement unmarshalling of the NullInt type
func (nullInt *NullInt) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON int

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	nullInt.Value = unmarshalledJSON
	nullInt.HasValue = true

	return nil
}
