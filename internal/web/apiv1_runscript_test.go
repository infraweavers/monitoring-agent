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

type RunScriptTestCase struct {
	ScriptToRun
	ExpectedResult
}

type RunScriptWithTimeoutTestCase struct {
	ScriptToRun
	Timeout
}

var osSpecifiTestCases = map[string]RunScriptTestCase{
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

func TestRunscriptApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {

		assert := assert.New(t)

		jsonBody, _ := json.Marshal(osSpecifiTestCases[runtime.GOOS].ScriptToRun)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecifiTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody, "Response output did not match case")
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be BadRequest")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 400 bad request with stdin supplied", func(t *testing.T) {
		var scriptBodyToTest = ScriptAsStdInToRun{
			ScriptToRun{
				Path: "sh",
				Args: []string{"sh", "-s"},
			},
			StdIn{
				StdIn: "uname",
			},
			StdInSignature{
				StdInSignature: "",
			},
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(scriptBodyToTest)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - This endpoint does not use stdin"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 400 bad request with invalid supplied", func(t *testing.T) {
		var scriptBodyToTest = RunScriptWithTimeoutTestCase{
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
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Invalid timeout supplied: '2'"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 200 and \"exitcode: 3, output: The script timed out\" when timeout exceeded", func(t *testing.T) {
		var osSpecifiTestCases = map[string]RunScriptWithTimeoutTestCase{
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
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be Service Unavailable")
		assert.Equal(`{"exitcode":3,"output":"The script timed out"}`, output.ResponseBody)
	})

	t.Run("Returns HTTP status 200 when timeout not exceeded", func(t *testing.T) {
		var osSpecifiTestCases = map[string]RunScriptWithTimeoutTestCase{
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
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":""}`, output.ResponseBody)
	})

	t.Run("200 Response to Approved Path and Arguments", func(t *testing.T) {
		configuration.Settings.Security.ApprovedPathArgumentsOnly = true
		configuration.Settings.Security.ApprovedPathArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", `write-host "Hello, World"`}, {"-command"}},
			"sh": {{"-c", "uname"}},
		}

		osSpecificRunScript := osSpecifiTestCases[runtime.GOOS].ScriptToRun

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(osSpecifiTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody)
	})

	t.Run("Bad request due to invalid path/arg combo", func(t *testing.T) {
		configuration.Settings.Security.ApprovedPathArgumentsOnly = true
		configuration.Settings.Security.ApprovedPathArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", "start-sleep 1"}, {"-command"}},
			"sh": {{"-c", "-s"}},
		}

		osSpecificRunScript := osSpecifiTestCases[runtime.GOOS].ScriptToRun

		jsonBody, _ := json.Marshal(osSpecificRunScript)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Unapproved Path/Args"}`, output.ResponseBody)
	})
}
