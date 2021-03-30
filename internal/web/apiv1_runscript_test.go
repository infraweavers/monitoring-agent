package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"testing"
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
		jsonBody, err := json.Marshal(osSpecificRunScriptTestCases[runtime.GOOS].ScriptToRun)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		if output.ResponseStatus != http.StatusOK {
			t.Errorf("Test Failed: Expected: %d, Got: %d", http.StatusOK, output.ResponseStatus)
		}

		if output.ResponseBody != osSpecificRunScriptTestCases[runtime.GOOS].ExpectedResult.Output {
			t.Errorf("Test Failed: Expected: %s, Got: %s", osSpecificRunScriptTestCases[runtime.GOOS].ExpectedResult.Output, output.ResponseBody)
		}
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		jsonBody, err := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		expectedResponseStatus := http.StatusBadRequest
		expectedResponseBody := `{"exitcode":3,"output":"400 Bad Request"}`
		output := TestHTTPRequestWithDefaultCredentials(t, request)

		if output.ResponseStatus != expectedResponseStatus {
			t.Errorf("Test Failed: Expected: %d, Got: %d", expectedResponseStatus, output.ResponseStatus)
		}

		if output.ResponseBody != expectedResponseBody {
			t.Errorf("Test Failed: Expected: %s, Got: %s", expectedResponseBody, output.ResponseBody)
		}
	})

	t.Run("Returns HTTP status 400 bad request with stdin supplied", func(t *testing.T) {
		var test = ScriptAsStdInToRun{
			ScriptToRun{
				Path: "sh",
				Args: []string{"sh", "-s"},
			},
			StdIn{
				StdIn: "uname",
			},
		}

		jsonBody, err := json.Marshal(test)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscript", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		expectedResponseStatus := http.StatusBadRequest
		expectedResponseBody := `{"exitcode":3,"output":"400 Bad Request - This endpoint does not use stdin"}`
		output := TestHTTPRequestWithDefaultCredentials(t, request)

		if output.ResponseStatus != expectedResponseStatus {
			t.Errorf("Test Failed: Expected: %d, Got: %d", expectedResponseStatus, output.ResponseStatus)
		}

		if output.ResponseBody != expectedResponseBody {
			t.Errorf("Test Failed: Expected: %s, Got: %s", expectedResponseBody, output.ResponseBody)
		}
	})
}
