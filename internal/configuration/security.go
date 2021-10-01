package configuration

import (
	"encoding/json"
)

// Security is a struct for representing security configuration, unmarshalled from the configuration.json file
type Security struct {
	DisableHTTPs                NullBool              `json:"DisableHTTPs" mandatory:"true"`
	SignedStdInOnly             NullBool              `json:"SignedStdInOnly" mandatory:"true"`
	MiniSign                    MiniSign              `json:"PublicKey" mandatory:"true"`
	AllowedAddresses            AllowedNetworks       `json:"AllowedAddresses" mandatory:"true"`
	UseClientCertificates       NullBool              `json:"UseClientCertificates" mandatory:"true"`
	ClientCertificateCAFile     ClientCertCA          `json:"ClientCertificateCAFile" mandatory:"true"`
	ApprovedExecutablesOnly     NullBool              `json:"ApprovedExecutablesOnly" mandatory:"true"`
	ApprovedExecutableArguments map[string][][]string `json:"ApprovedExecutableArguments" mandatory:"true"`
	AllowScriptArguments        NullBool              `json:"AllowScriptArguments" mandatory:"true"`
}

// UnmarshalJSON is a method to implement unmarshalling of the Security type
func (security *Security) UnmarshalJSON(b []byte) error {
	type JSONTmp Security
	var jsonTmp JSONTmp

	err := json.Unmarshal(b, &jsonTmp)
	if err != nil {
		return err
	}

	*security = Security(jsonTmp)
	err = validateStruct(security)
	if err != nil {
		return err
	}

	return nil
}
