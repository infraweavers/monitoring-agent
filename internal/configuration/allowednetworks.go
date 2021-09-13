package configuration

import (
	"encoding/json"
	"net"
)

type AllowedNetworks struct {
	CIDR []*net.IPNet
}

func (allowedNetworks *AllowedNetworks) UnmarshalJSON(b []byte) error {
	var unmarshalledJson []string

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	for x := 0; x < len(unmarshalledJson); x++ {
		_, network, err := net.ParseCIDR(unmarshalledJson[x])
		if err != nil {
			return err
		}
		allowedNetworks.CIDR = append(allowedNetworks.CIDR, network)
	}
	return nil
}
