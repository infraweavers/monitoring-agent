package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"testing"
)

type ScriptAsStdInToRun struct {
	Path  string
	Args  []string
	StdIn string
}

type RunScriptStdInTestCase struct {
	Path     string
	Args     []string
	StdIn    string
	Expected string
}

var runScriptStdinTestCases = map[string]RunScriptStdInTestCase{
	"linux": {
		Path:     "sh",
		Args:     []string{"-s"},
		StdIn:    "uname",
		Expected: `{"exitcode":0,"output":"Linux\n"}`,
	},
	"windows": {
		Path:     `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
		Args:     []string{"-command", "-"},
		StdIn:    `Write-Host "Hello-World"`,
		Expected: `{"exitcode":0,"output":"Hello, World\n"}`,
	},
}

func TestRunScriptStdInApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		testCase := runScriptStdinTestCases[runtime.GOOS]
		var testScript = ScriptAsStdInToRun{
			Path:  testCase.Path,
			Args:  testCase.Args,
			StdIn: testCase.StdIn,
		}

		jsonBody, err := json.Marshal(testScript)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		output := TestHTTPRequestWithDefaultCredentials(t, request)

		if output.ResponseStatus != http.StatusOK {
			t.Errorf("Test Failed: Expected: %d, Got: %d", http.StatusOK, output.ResponseStatus)
		}

		if output.ResponseBody != testCase.Expected {
			t.Errorf("Test Failed: Expected: %s, Got: %s", testCase.Expected, output.ResponseBody)
		}
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		jsonBody, err := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, GetTestServerURL(t)+"/v1/runscriptstdin", bytesBuffer)
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
}
