package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("A valid minisign signature should be treated as signed", func(t *testing.T) {
		testBody := "Write-Host 'Hello, World'"
		testSignature := "untrusted comment: signature from minisign secret key\nRWQ3ly9IPenQ6XE4gvV0tpJPSRdw/Si+Q4r97LbpLj0Hb3sV+XFydynJg3iFT2PjIlE3xViNOmFT9XrIoidedDr41+Ly0AYbUQg=\ntrusted comment: timestamp:1617721023	file:robtest.ps1\nHkxuqHSvipJo/unNKgDS+JGDB0+Q5d8nOeoJ0NGOnKBNsNdvAj8FWf7fhaPV7mzRJ1ooLvYpI0yUsD7lpaDwBQ==\n"

		response := verifySignature(testBody, testSignature)
		assert.True(t, response)
	})

	t.Run("A valid minisign signature should be treated as signed", func(t *testing.T) {
		testBody := "Write-Host 'Hello, World' Altered String"
		testSignature := "untrusted comment: signature from minisign secret key\nRWQ3ly9IPenQ6XE4gvV0tpJPSRdw/Si+Q4r97LbpLj0Hb3sV+XFydynJg3iFT2PjIlE3xViNOmFT9XrIoidedDr41+Ly0AYbUQg=\ntrusted comment: timestamp:1617721023	file:robtest.ps1\nHkxuqHSvipJo/unNKgDS+JGDB0+Q5d8nOeoJ0NGOnKBNsNdvAj8FWf7fhaPV7mzRJ1ooLvYpI0yUsD7lpaDwBQ==\n"

		response := verifySignature(testBody, testSignature)
		assert.False(t, response)
	})
}
