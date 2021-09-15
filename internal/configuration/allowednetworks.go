package configuration

import (
	"encoding/json"
	"net"
)

// AllowedNetworks is a struct representing a list of CIDRs permitted to use monitoring-agent, unmarshalled from the configuration.json file
type AllowedNetworks struct {
	CIDR []*net.IPNet
}

// UnmarshalJSON is a method to implement unmarshalling of the AllowedNetworks type
func (allowedNetworks *AllowedNetworks) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON []string

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	for x := 0; x < len(unmarshalledJSON); x++ {
		_, network, err := net.ParseCIDR(unmarshalledJSON[x])
		if err != nil {
			return err
		}
		allowedNetworks.CIDR = append(allowedNetworks.CIDR, network)
	}
	return nil
}
