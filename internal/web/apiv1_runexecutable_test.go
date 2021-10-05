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

type RunExecutableTestCase struct {
	ScriptToRun
	ExpectedResult
}

type RunExecutableWithTimeoutTestCase struct {
	ScriptToRun
	Timeout
}

var osSpecifiTestCases = map[string]RunExecutableTestCase{
	"linux": {
		ScriptToRun{
			Path: "sh",
			Args: []string{"-c", "uname"},
		},
		ExpectedResult{
			Output: `{"exitcode":0,"output":"Linux\n"}`,
		},
	},
	"windows": {
		ScriptToRun{
			Path: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			Args: []string{"-command", `write-host "Hello, World"`},
		},
		ExpectedResult{
			Output: `{"exitcode":0,"output":"Hello, World\n"}`,
		},
	},
}

func TestRunexecutableApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {

		assert := assert.New(t)

		jsonBody, _ := json.Marshal(osSpecifiTestCases[runtime.GOOS].ScriptToRun)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecifiTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody, "Response output did not match case")
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be BadRequest")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = false
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}
		expectedOutput := ""

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `Write-Host 'Hello, World'`,
				"stdinsignature": "untrusted comment: signature from minisign secret key\nRWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=\ntrusted comment: timestamp:1629361915	file:writehost.txt\nOfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==\n",
			}
			expectedOutput = `{"exitcode":0,"output":"Hello, World\n"}`
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `uname`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI8mVzlQxqbNt9+ldPNoPREsedr+sAHAnkrkyg80yQo1UrrYD7+ScU9ZXqYv79ukLN3nEgK8tsQ4uUSH7Sgpw1AY=
trusted comment: timestamp:1629361789	file:uname.txt
6ZxQL0d64hC8LCCPpKct+oyPN/JV1zqnD+92Uk9z9dEYnugpYmgVv9ZXabaLePEIP3bfNYe5JeD83YHWYS4/Aw==
`,
			}
			expectedOutput = `{"exitcode":0,"output":"Linux\n"}`
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(expectedOutput, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Returns HTTP status 400 bad request with stdin supplied", func(t *testing.T) {
		scriptBodyToTest := map[string]interface{}{
			"path":           `sh`,
			"args":           []string{"-s"},
			"stdin":          `uname`,
			"stdinsignature": ``,
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(scriptBodyToTest)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 400 bad request with invalid supplied", func(t *testing.T) {
		var scriptBodyToTest = RunExecutableWithTimeoutTestCase{
			ScriptToRun{
				Path: "sh",
				Args: []string{"sh", "-s"},
			},
			Timeout{
				Timeout: "2",
			},
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(scriptBodyToTest)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Invalid timeout supplied: '2'"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 200 and \"exitcode: 3, output: The script timed out\" when timeout exceeded", func(t *testing.T) {
		var osSpecifiTestCases = map[string]RunExecutableWithTimeoutTestCase{
			"linux": {
				ScriptToRun{
					Path: "sh",
					Args: []string{"-c", "sleep 2"},
				},
				Timeout{
					Timeout: "1s",
				},
			},
			"windows": {
				ScriptToRun{
					Path: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
					Args: []string{"-command", `Start-Sleep -seconds 2`},
				},
				Timeout{
					Timeout: "1s",
				},
			},
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(osSpecifiTestCases[runtime.GOOS])
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be Service Unavailable")
		assert.Equal(`{"exitcode":3,"output":"The script timed out"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 200 when timeout not exceeded", func(t *testing.T) {
		var osSpecifiTestCases = map[string]RunExecutableWithTimeoutTestCase{
			"linux": {
				ScriptToRun{
					Path: "sh",
					Args: []string{"-c", "sleep 1"},
				},
				Timeout{
					Timeout: "2s",
				},
			},
			"windows": {
				ScriptToRun{
					Path: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
					Args: []string{"-command", `Start-Sleep -seconds 1`},
				},
				Timeout{
					Timeout: "2s",
				},
			},
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(osSpecifiTestCases[runtime.GOOS])
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":""}`, output.ResponseBody)
	})

	t.Run("200 Response to Approved Path and Arguments", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.ApprovedExecutableArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", `write-host "Hello, World"`}, {"-command"}},
			"sh": {{"-c", "uname"}},
		}

		osSpecificRunExecutable := osSpecifiTestCases[runtime.GOOS].ScriptToRun

		jsonBody, _ := json.Marshal(osSpecificRunExecutable)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecifiTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody)
	})

	t.Run("Bad request due to invalid path/arg combo", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.ApprovedExecutableArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", "start-sleep 1"}, {"-command"}},
			"sh": {{"-c", "-s"}},
		}

		osSpecificRunExecutable := osSpecifiTestCases[runtime.GOOS].ScriptToRun

		jsonBody, _ := json.Marshal(osSpecificRunExecutable)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runexecutable", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Unapproved Path/Args"}`, output.ResponseBody)
	})
}
