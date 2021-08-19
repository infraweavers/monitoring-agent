package web

import (
	"monitoringagent/internal/configuration"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("A valid minisign signature should be treated as signed", func(t *testing.T) {
		configuration.TestingInitialise()

		testBody := `Write-Host 'Hello, World'`
		testSignature := `untrusted comment: signature from minisign secret key
RWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=
trusted comment: timestamp:1629361915	file:writehost.txt
OfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==
`

		response := verifySignature(testBody, testSignature)
		assert.True(t, response)
	})

	t.Run("A valid minisign signature should be treated as signed", func(t *testing.T) {
		configuration.TestingInitialise()

		testBody := `Write-Host 'Hello, World' Altered String`
		testSignature := `untrusted comment: signature from minisign secret key
RWQ3ly9IPenQ6XE4gvV0tpJPSRdw/Si+Q4r97LbpLj0Hb3sV+XFydynJg3iFT2PjIlE3xViNOmFT9XrIoidedDr41+Ly0AYbUQg=
trusted comment: timestamp:1617721023	file:robtest.ps1
HkxuqHSvipJo/unNKgDS+JGDB0+Q5d8nOeoJ0NGOnKBNsNdvAj8FWf7fhaPV7mzRJ1ooLvYpI0yUsD7lpaDwBQ==
`

		response := verifySignature(testBody, testSignature)
		assert.False(t, response)
	})
}

func TestVerifyRemoteHost(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("Empty AllowedAddresses should block all IPs", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{}

		assert.False(t, verifyRemoteHost("127.0.0.1:5431"))
	})

	t.Run("Loopback should be permitted when granted", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{
			{IP: net.ParseIP("127.0.0.0"), Mask: net.IPMask(net.ParseIP("255.0.0.0").To4())},
		}
		assert.True(t, verifyRemoteHost("127.0.0.1:5431"))
	})

	t.Run("0.0.0.0/0 should grant all", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{
			{IP: net.ParseIP("0.0.0.0"), Mask: net.IPMask(net.ParseIP("0.0.0.0").To4())},
		}
		assert.True(t, verifyRemoteHost("54.74.60.0:6428"))
	})

	t.Run("An IP outside of the range should be denied", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{
			{IP: net.ParseIP("192.168.54.0"), Mask: net.IPMask(net.ParseIP("255.255.255.0").To4())},
		}
		assert.False(t, verifyRemoteHost("8.8.8.8:9945"))
	})

	t.Run("Only IPs within any of the masks should be granted ", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{
			{IP: net.ParseIP("192.168.54.0"), Mask: net.IPMask(net.ParseIP("255.255.255.0").To4())},
			{IP: net.ParseIP("192.168.51.0"), Mask: net.IPMask(net.ParseIP("255.255.255.0").To4())},
			{IP: net.ParseIP("8.8.0.0"), Mask: net.IPMask(net.ParseIP("255.255.0.0").To4())},
		}
		assert.True(t, verifyRemoteHost("8.8.8.8:9945"))
		assert.True(t, verifyRemoteHost("192.168.51.54:1"))
		assert.True(t, verifyRemoteHost("192.168.54.204:1"))

		assert.False(t, verifyRemoteHost("127.0.0.1:5421"))
		assert.False(t, verifyRemoteHost("10.1.1.1:5874"))
		assert.False(t, verifyRemoteHost("[::1]:5421"))
		assert.False(t, verifyRemoteHost("[ec9b:434a:0623:2620:9fa3:5432:ee23:ea81]:5421"))
	})

	t.Run("IPv6 should be supported", func(t *testing.T) {
		configuration.TestingInitialise()
		configuration.Settings.AllowedAddresses = []*net.IPNet{
			{IP: net.ParseIP("192.168.54.0"), Mask: net.IPMask(net.ParseIP("255.255.255.0").To4())},
			{IP: net.ParseIP("ec9b:434a:0623:2620::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::").To16())},
			{IP: net.ParseIP("127.0.0.0"), Mask: net.IPMask(net.ParseIP("255.0.0.0").To4())},
		}

		assert.False(t, verifyRemoteHost("[::1]:5421"))
		assert.True(t, verifyRemoteHost("[ec9b:434a:0623:2620:9fa3:5432:ee23:ea81]:5421"))
	})

}
