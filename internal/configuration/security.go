package configuration

import (
	"encoding/json"
)

// JSONconfigSecurity is a struct for unmarshalling the configuration.json file
type Security struct {
	DisableHTTPs              NullBool              `json:"DisableHTTPs" mandatory:"true"`
	SignedStdInOnly           NullBool              `json:"SignedStdInOnly" mandatory:"true"`
	MiniSign                  MiniSign              `json:"PublicKey" mandatory:"true"`
	AllowedAddresses          AllowedNetworks       `json:"AllowedAddresses" mandatory:"true"`
	UseClientCertificates     NullBool              `json:"UseClientCertificates" mandatory:"true"`
	ClientCertificateCAFile   ClientCertCA          `json:"ClientCertificateCAFile" mandatory:"true"`
	ApprovedPathArgumentsOnly NullBool              `json:"ApprovedPathArgumentsOnly" mandatory:"true"`
	ApprovedPathArguments     map[string][][]string `json:"ApprovedPathArguments" mandatory:"true"`
}

func (security *Security) UnmarshalJSON(b []byte) error {
	type JsonTmp Security
	var jsonTmp JsonTmp

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
