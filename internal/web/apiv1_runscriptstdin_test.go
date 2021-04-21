package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ScriptAsStdInToRun struct {
	ScriptToRun
	StdIn
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
				StdIn: `Write-Host "Hello, World"`,
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
}
