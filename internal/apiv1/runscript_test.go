package apiv1

import (
	"bytes"
	"encoding/json"
	"mama/internal/testhelpers"
	"net/http"
	"runtime"
	"testing"
)

type ScriptToRun struct {
	Path string
	Args []string
}

type TestCase struct {
	Path     string
	Args     []string
	Expected string
}

var testCases = map[string]TestCase{
	"Linux": {
		Path:     "/bin/sh",
		Args:     []string{"-c", "/usr/bin/uname"},
		Expected: `{"exitcode":0,"output":"Linux\n"}`,
	},
	"windows": {
		Path:     `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
		Args:     []string{"-command", `write-host "Hello, World"`},
		Expected: `{"exitcode":0,"output":"Hello, World\n"}`,
	},
}

func TestRunscriptApiHandler(t *testing.T) {

	testhelpers.Setup(RunscriptPostHandler)
	defer testhelpers.Teardown()

	t.Run("Runs supplied script, returns HTTP status 200 and expected script output", func(t *testing.T) {
		testCase := testCases[runtime.GOOS]
		var testScript = ScriptToRun{
			Path: testCase.Path,
			Args: testCase.Args,
		}

		jsonBody, err := json.Marshal(testScript)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, testhelpers.GetServerURL(t)+"/v1/runscript/", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		output := testhelpers.TestHTTPRequest(t, request)

		if output.ResponseStatus != http.StatusOK {
			t.Error("Test Failed: Expected: {}, Got: {}", http.StatusOK, output.ResponseStatus)
		}

		if output.ResponseBody != testCase.Expected {
			t.Error("Test Failed: Expected: {}, Got: {}", testCase.Expected, output.ResponseBody)
		}
	})

	t.Run("Returns HTTP status 400 bad request with erronous post", func(t *testing.T) {
		jsonBody, err := json.Marshal(`{"foo": "bar", "doh": "car"}`)
		if err != nil {
			t.Fatal(err)
		}

		bytesBuffer := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, testhelpers.GetServerURL(t)+"/v1/runscript/", bytesBuffer)
		if err != nil {
			t.Fatal(err)
		}

		expectedResponseStatus := http.StatusBadRequest
		expectedResponseBody := `{"exitcode":3,"output":"400 Bad Request"}`
		output := testhelpers.TestHTTPRequest(t, request)

		if output.ResponseStatus != expectedResponseStatus {
			t.Error("Test Failed: Expected: {}, Got: {}", expectedResponseStatus, output.ResponseStatus)
		}

		if output.ResponseBody != expectedResponseBody {
			t.Error("Test Failed: Expected: {}, Got: {}", expectedResponseBody, output.ResponseBody)
		}
	})
}

// ToDo: TestRunscriptApiHandlerWithEmptyBody
