package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RunScriptTestCase struct {
	ScriptToRun
	ExpectedResult
}

var osSpecificRunScriptTestCases = map[string]RunScriptTestCase{
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

		jsonBody, _ := json.Marshal(osSpecificRunScriptTestCases[runtime.GOOS].ScriptToRun)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(output.ResponseStatus, http.StatusOK, "Response code should be OK")
		assert.Equal(output.ResponseBody, osSpecificRunScriptTestCases[runtime.GOOS].ExpectedResult.Output, "Response output did not match case")
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(output.ResponseStatus, http.StatusBadRequest, "Response code should be BadRequest")
		assert.Equal(output.ResponseBody, `{"exitcode":3,"output":"400 Bad Request"}`)
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
		}
		assert := assert.New(t)

		jsonBody, _ := json.Marshal(scriptBodyToTest)
		request, _ := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytes.NewBuffer(jsonBody))

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		assert.Equal(output.ResponseStatus, http.StatusBadRequest, "Response code should be Bad Request")
		assert.Equal(output.ResponseBody, `{"exitcode":3,"output":"400 Bad Request - This endpoint does not use stdin"}`)
	})
}
