package web

import (
	"monitoringagent/internal/configuration"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunexecutableApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("When supplied an executable, returns HTTP status 200 and expected script output", func(t *testing.T) {
		testRequest := map[string]interface{}{}
		expectedResult := ""

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args": []string{"-command", `write-host "Hello, World"`},
			}
			expectedResult = `{"exitcode":0,"output":"Hello, World\n"}`
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path": `sh`,
				"args": []string{"-c", "uname"},
			}
			expectedResult = `{"exitcode":0,"output":"Linux\n"}`
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(expectedResult, output.ResponseBody, "Response output did not match case")
	})

	t.Run("When provided with a entirely incorrect request we return HTTP status 400 bad request", func(t *testing.T) {

		testRequest := map[string]interface{}{
			"foo": `bar`,
			"doh": `car`,
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be BadRequest")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("When provided a request with stdin an  HTTP status 400 bad request is returned", func(t *testing.T) {

		testRequest := map[string]interface{}{
			"path":           `sh`,
			"args":           []string{"-s"},
			"stdin":          `uname`,
			"stdinsignature": ``,
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("When provided with an invalid timeout, returns HTTP status 400 bad request", func(t *testing.T) {

		testRequest := map[string]interface{}{
			"path":    `sh`,
			"args":    []string{"-c", "-s"},
			"Timeout": "2",
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Invalid timeout supplied: '2'"}`, output.ResponseBody)
	})

	t.Run(`When the underlying command times out an HTTP status 200 and "exitcode: 3, output: The script timed out" is returned`, func(t *testing.T) {

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":    `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":    []string{"-command", `Start-Sleep -seconds 2`},
				"Timeout": "1s",
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":    `sh`,
				"args":    []string{"-c", "sleep 2"},
				"Timeout": "1s",
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be Service Unavailable")
		assert.Equal(`{"exitcode":3,"output":"The script timed out"}`, output.ResponseBody)
	})

	t.Run("When the command completes with in the timeout an HTTP 200 OK is returned", func(t *testing.T) {

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":    `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":    []string{"-command", `Start-Sleep -seconds 1`},
				"Timeout": "2s",
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":    `sh`,
				"args":    []string{"-c", "sleep 1"},
				"Timeout": "2s",
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":""}`, output.ResponseBody)
	})

	t.Run("When the submitted request is included in the ApprovedExecutables and the ApprovedExecutablesOnly is enabled then a HTTP 200 OK Response with the command output is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.ApprovedExecutableArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", `write-host "Hello, World"`}, {"-command"}},
			"sh": {{"-c", "uname"}},
		}

		testRequest := map[string]interface{}{}
		expectedOutput := ""

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args": []string{"-command", `write-host "Hello, World"`},
			}
			expectedOutput = `{"exitcode":0,"output":"Hello, World\n"}`
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path": `sh`,
				"args": []string{"-c", "uname"},
			}
			expectedOutput = `{"exitcode":0,"output":"Linux\n"}`
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(expectedOutput, output.ResponseBody)
	})

	t.Run("When the submitted request is not included in the ApprovedExecutables and the ApprovedExecutablesOnly is enabled then the request is rejected", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.ApprovedExecutableArguments = map[string][][]string{
			`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`: {{"-command", "-"}, {"-command", "start-sleep 1"}, {"-command"}},
			"sh": {{"-c", "-s"}},
		}

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args": []string{"-command", `write-host "Hello, World"`},
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path": `sh`,
				"args": []string{"-c", "uname"},
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runexecutable", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Unapproved Path/Args"}`, output.ResponseBody)
	})

}
