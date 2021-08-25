package web

import (
	"bytes"
	"encoding/json"
	"monitoringagent/internal/configuration"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ScriptAsStdInToRun struct {
	ScriptToRun
	StdIn
	StdInSignature
}

type RunScriptStdInTestCase struct {
	ScriptAsStdInToRun
	ExpectedResult
}

var osSpecificRunScriptStdinTestCases = map[string]RunScriptStdInTestCase{
	"linux": {
		ScriptAsStdInToRun{
			ScriptToRun{
				Path: "sh",
				Args: []string{"-s"},
			},
			StdIn{
				StdIn: "uname",
			},
			StdInSignature{
				StdInSignature: `untrusted comment: signature from minisign secret key
RWTV8L06+shYI8mVzlQxqbNt9+ldPNoPREsedr+sAHAnkrkyg80yQo1UrrYD7+ScU9ZXqYv79ukLN3nEgK8tsQ4uUSH7Sgpw1AY=
trusted comment: timestamp:1629361789	file:uname.txt
6ZxQL0d64hC8LCCPpKct+oyPN/JV1zqnD+92Uk9z9dEYnugpYmgVv9ZXabaLePEIP3bfNYe5JeD83YHWYS4/Aw==
`,
			},
		},
		ExpectedResult{
			Output: `{"exitcode":0,"output":"Linux\n"}`,
		},
	},
	"windows": {
		ScriptAsStdInToRun{
			ScriptToRun{
				Path: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				Args: []string{"-command", "-"},
			},
			StdIn{
				StdIn: `Write-Host 'Hello, World'`,
			},
			StdInSignature{
				StdInSignature: `untrusted comment: signature from minisign secret key
RWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=
trusted comment: timestamp:1629361915	file:writehost.txt
OfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==
`,
			},
		},
		ExpectedResult{
			Output: `{"exitcode":0,"output":"Hello, World\n"}`,
		},
	},
}

func TestRunScriptStdInApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		jsonBody, _ := json.Marshal(osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecificRunScriptStdinTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		jsonBody, _ := json.Marshal(`{"foo": "bar", "doh": "car"}`)

		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))
		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 400 bad request without stdin supplied", func(t *testing.T) {
		var test = ScriptToRun{
			Path: "sh",
			Args: []string{"sh", "-s"},
		}

		jsonBody, _ := json.Marshal(test)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - This endpoint requires stdin"}`, output.ResponseBody)
	})

	t.Run("Runs supplied signed script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		configuration.Settings.SignedStdInOnly = true

		jsonBody, _ := json.Marshal(osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecificRunScriptStdinTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Runs unsigned script, returns HTTP status 400 and expected failed signiture verification", func(t *testing.T) {
		configuration.Settings.SignedStdInOnly = true

		osSpecificRunScript := osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun
		if runtime.GOOS == "windows" {
			osSpecificRunScript.StdIn.StdIn = `Write-Host 'This script is a test.'`
		}
		if runtime.GOOS == "linux" {
			osSpecificRunScript.StdIn.StdIn = `echo This script is a test."`
		}

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Signature not valid"}`, output.ResponseBody)
	})

	t.Run("Runs supplied signed script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		configuration.Settings.SignedStdInOnly = true

		osSpecificRunScript := osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun
		if runtime.GOOS == "windows" {
			osSpecificRunScript.StdIn.StdIn = `Write-Host 'This script is a test.'`
			osSpecificRunScript.StdInSignature.StdInSignature = `untrusted comment: signature from minisign secret key
RWTV8L06+shYI0gq2Ph8MRbdPBrxVEXwzw12yn6b6qG4uyBcnCZ6jTBVULVTZPlMwx6mBnLL2ayCwL/NC83wHJMBtcg3oY/uDQk=
trusted comment: timestamp:1629362484	file:whtest.txt
tbOXpkm9GyEQlUflmVX4cDy2k5fJWU3wtxscvAqSu19C227SFQU6SHlUZbpXB85pBoFJTJK+tQVBN1u1RmaOCw==
`
		}
		if runtime.GOOS == "linux" {
			osSpecificRunScript.StdIn.StdIn = `echo "This script is a test."`
			osSpecificRunScript.StdInSignature.StdInSignature = `untrusted comment: signature from minisign secret key
RWTV8L06+shYI+aL2MAm12HN97gM83Cd1c2H10PMtGhFAmYlxsEWnJGZEFMyFtB46Ity/6iK36IEw66L+5KjcLJEOhw7TMwjZQs=
trusted comment: timestamp:1629368840	file:echotest.txt
U54CjtRd9nA/jp4iEhdbQ35eE4yWQRY0nbJlw4elRwilslde8nrZwfaIK1a2R+7gzfeuiZq8xTlKtIvTOg5aAA==
`
		}

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":"This script is a test.\n"}`, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Runs supplied signed script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		configuration.Settings.ApprovedPathArgumentsOnly = true
		configuration.Settings.SignedStdInOnly = true

		osSpecificRunScript := osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecificRunScriptStdinTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Bad request due to invalid path/arg combo", func(t *testing.T) {
		configuration.Settings.ApprovedPathArgumentsOnly = true
		configuration.Settings.SignedStdInOnly = true

		osSpecificRunScript := osSpecificRunScriptStdinTestCases[runtime.GOOS].ScriptAsStdInToRun
		osSpecificRunScript.ScriptToRun.Args = append(osSpecificRunScript.ScriptToRun.Args, "")

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Unapproved Path/Args"}`, output.ResponseBody)
	})
}
