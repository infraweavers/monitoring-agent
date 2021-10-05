package web

import (
	"monitoringagent/internal/configuration"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunScriptStdInApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("When script is supplied through stdin, returns HTTP status 200 and expected script output", func(t *testing.T) {
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
			}
			expectedOutput = `{"exitcode":0,"output":"Hello, World\n"}`
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `uname`,
			}
			expectedOutput = `{"exitcode":0,"output":"Linux\n"}`
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(expectedOutput, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("When provided a request with totally the wrong shape, a 400 is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = false

		testRequest := map[string]interface{}{
			"foo": "bar",
			"doh": "car",
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request"}`, output.ResponseBody)
	})

	t.Run("When provided a request missing just stdin, a 400 is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args": []string{"-command", "-"},
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path": `sh`,
				"args": []string{"-s"},
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Only signed stdin can be executed"}`, output.ResponseBody)
	})

	t.Run("When provided a stdin and a valid signature, the script should be executed, an HTTP 200 returned, exit code of 0 and the output in the response", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
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

	t.Run("When provided with stdin but an invalid signature, returns HTTP status 400 and expected failed signiture verification", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `Write-Host 'This script is an incorrectly signed test.'`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI0gq2Ph8MRbdPBrxVEXwzw12yn6b6qG4uyBcnCZ6jTBVULVTZPlMwx6mBnLL2ayCwL/NC83wHJMBtcg3oY/uDQk=
trusted comment: timestamp:1629362484	file:whtest.txt
tbOXpkm9GyEQlUflmVX4cDy2k5fJWU3wtxscvAqSu19C227SFQU6SHlUZbpXB85pBoFJTJK+tQVBN1u1RmaOCw==
`,
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `echo "This script is an incorrectly signed test."`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI+aL2MAm12HN97gM83Cd1c2H10PMtGhFAmYlxsEWnJGZEFMyFtB46Ity/6iK36IEw66L+5KjcLJEOhw7TMwjZQs=
trusted comment: timestamp:1629368840	file:echotest.txt
U54CjtRd9nA/jp4iEhdbQ35eE4yWQRY0nbJlw4elRwilslde8nrZwfaIK1a2R+7gzfeuiZq8xTlKtIvTOg5aAA==
`,
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Signature not valid"}`, output.ResponseBody)
	})

	t.Run("When provided provided stdin and a valid signature the supplied signed script is executed, an HTTP 200 and the output from the script are returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `Write-Host 'This script is a test.'`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI0gq2Ph8MRbdPBrxVEXwzw12yn6b6qG4uyBcnCZ6jTBVULVTZPlMwx6mBnLL2ayCwL/NC83wHJMBtcg3oY/uDQk=
trusted comment: timestamp:1629362484	file:whtest.txt
tbOXpkm9GyEQlUflmVX4cDy2k5fJWU3wtxscvAqSu19C227SFQU6SHlUZbpXB85pBoFJTJK+tQVBN1u1RmaOCw==
`,
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `echo "This script is a test."`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI+aL2MAm12HN97gM83Cd1c2H10PMtGhFAmYlxsEWnJGZEFMyFtB46Ity/6iK36IEw66L+5KjcLJEOhw7TMwjZQs=
trusted comment: timestamp:1629368840	file:echotest.txt
U54CjtRd9nA/jp4iEhdbQ35eE4yWQRY0nbJlw4elRwilslde8nrZwfaIK1a2R+7gzfeuiZq8xTlKtIvTOg5aAA==
`,
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":"This script is a test.\n"}`, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("When approved executables are enabled but stdin signatures are disabled and provided with a whitelisted command, a success response is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = false
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":           `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":           []string{"-command", "-"},
				"stdin":          `Write-Host 'This script is a test.'`,
				"stdinsignature": ``,
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":           `sh`,
				"args":           []string{"-s"},
				"stdin":          `echo "This script is a test."`,
				"stdinsignature": ``,
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":"This script is a test.\n"}`, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("When approved executables are enabled and stdin signatures are enabled and provided with a whitelisted command, a success response is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `Write-Host 'This script is a test.'`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI0gq2Ph8MRbdPBrxVEXwzw12yn6b6qG4uyBcnCZ6jTBVULVTZPlMwx6mBnLL2ayCwL/NC83wHJMBtcg3oY/uDQk=
trusted comment: timestamp:1629362484	file:whtest.txt
tbOXpkm9GyEQlUflmVX4cDy2k5fJWU3wtxscvAqSu19C227SFQU6SHlUZbpXB85pBoFJTJK+tQVBN1u1RmaOCw==
`,
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `echo "This script is a test."`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI+aL2MAm12HN97gM83Cd1c2H10PMtGhFAmYlxsEWnJGZEFMyFtB46Ity/6iK36IEw66L+5KjcLJEOhw7TMwjZQs=
trusted comment: timestamp:1629368840	file:echotest.txt
U54CjtRd9nA/jp4iEhdbQ35eE4yWQRY0nbJlw4elRwilslde8nrZwfaIK1a2R+7gzfeuiZq8xTlKtIvTOg5aAA==
`,
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(`{"exitcode":0,"output":"This script is a test.\n"}`, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("When approved executables are enabled and stdin signatures are enabled and provided with a whitelisted command, a success response and command output is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{}

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `(Get-CimInstance Win32_OperatingSystem).version`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI+dzaUD0xCAbUp0KBF9B+u5wiBaqe1ppiXLsVyyWAyfKXVo0q3pgLWvwkIvjTNk+q5OjrS6G4rclJU2mmP1v6wM=
trusted comment: timestamp:1629903111	file:winver.txt
XMMjgNkS+rnnAkC4gARhK1o83VB3pIAtOQAzO/RZ31x5HfgpWvZe0rAjO7hauH4mMBwjYL/71cqul4yPknnrAw==
`,
			}
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `uname -a`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI0/LRKnu1ask22XycnLwTEaCVyo3COcMqVJOYgi4VjkEYvNz6VLnWNzSqSqVNwCv6WkJwp6viFKBedcRKBfuGQ4=
trusted comment: timestamp:1629903102	file:uname.txt
JkeUlACQaVsrlHmFWg0U0Y5AcnbusFKHNF4bF3kGyixXS3B3/fCZ9T9LMyMbPwZyUJyMGBpfAVXgAQQdM82HCA==
`,
			}
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
	})

	t.Run("When approved executables and arguments and stdin signatures are enabled and provided with a whitelisted command, a success response and command output is returned", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
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

	t.Run("When ApprovedExecutablesOnly is enabled and we are provided with arguments that don't match then the request is rejected", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{
			"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			"args":  []string{"-wrong", "args"},
			"stdin": `Write-Host 'Hello, World'`,
			"stdinsignature": "untrusted comment: signature from minisign secret key\nRWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=\ntrusted comment: timestamp:1629361915	file:writehost.txt\nOfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==\n",
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus)
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Unapproved Path/Args"}`, output.ResponseBody)
	})

	t.Run("When ScriptArguments are disabled, requests with script arguments should be rejected", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{
			"path":            `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			"args":            []string{"-command", "-"},
			"ScriptArguments": []string{"scriptArgument"},
			"stdin":           `Write-Host 'Hello, World'`,
			"stdinsignature": "untrusted comment: signature from minisign secret key\nRWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=\ntrusted comment: timestamp:1629361915	file:writehost.txt\nOfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==\n",
		}
		expectedOutput := `{"exitcode":3,"output":"400 Bad Request - Script Arguments Passed But Not Permitted"}`

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(expectedOutput, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("When given an executable that does not exist, returns HTTP status 200 containing exit code 3 and the error output", func(t *testing.T) {
		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true

		testRequest := map[string]interface{}{}
		expectedOutput := ""

		if runtime.GOOS == "windows" {
			testRequest = map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `Write-Host 'Hello, World'`,
				"stdinsignature": "untrusted comment: signature from minisign secret key\nRWTV8L06+shYIx/hkk/yLgwyrJvVfYNoGDsCsv6/+2Tp1Feq/S6DLwpOENGpsUe15ZedtCZzjmXQrJ+vVeC2oNB3vR88G25o0wo=\ntrusted comment: timestamp:1629361915	file:writehost.txt\nOfDNTVG4KeQatDps8OzEXZGNhSQrfHOWTYJ2maNyrWe+TGss7VchEEFMrKMvvTP5q0NL9YoLvbyxoWxCd2H0Cg==\n",
			}
			expectedOutput = `{"exitcode":3,"output":"An error occurred executing the command: exec: \"sh\": executable file not found in %PATH%"}`
		}
		if runtime.GOOS == "linux" {
			testRequest = map[string]interface{}{
				"path":  `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
				"args":  []string{"-command", "-"},
				"stdin": `uname`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWTV8L06+shYI8mVzlQxqbNt9+ldPNoPREsedr+sAHAnkrkyg80yQo1UrrYD7+ScU9ZXqYv79ukLN3nEgK8tsQ4uUSH7Sgpw1AY=
trusted comment: timestamp:1629361789	file:uname.txt
6ZxQL0d64hC8LCCPpKct+oyPN/JV1zqnD+92Uk9z9dEYnugpYmgVv9ZXabaLePEIP3bfNYe5JeD83YHWYS4/Aw==
`,
			}
			expectedOutput = `{"exitcode":3,"output":"An error occurred executing the command: exec: \"C:\\\\Windows\\\\System32\\\\WindowsPowerShell\\\\v1.0\\\\powershell.exe\": executable file not found in $PATH"}`
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.Equal(expectedOutput, output.ResponseBody)
	})

	t.Run("When script arguments are disabled, passing scriptarguments should be rejected", func(t *testing.T) {

		configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = true
		configuration.Settings.Security.SignedStdInOnly.IsTrue = true
		configuration.Settings.Security.AllowScriptArguments.IsTrue = false

		testRequest := map[string]interface{}{
			"args":            []string{"-command", "-"},
			"path":            `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			"scriptarguments": []string{"arg1"},
			"stdin":           "#!/bin/bash\necho \"First line.\";\necho \"Second line.\";\necho \"Third lime.\";\n\npos=1\nfor arg in \"$@\"; do\n  echo \"$pos-th argument : $arg\"\n  (( pos += 1 ))\ndone\n\n\n",
			"timeout":         "10s",
		}

		output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusBadRequest, output.ResponseStatus, "Response code should be Bad Request")
		assert.Equal(`{"exitcode":3,"output":"400 Bad Request - Script Arguments Passed But Not Permitted"}`, output.ResponseBody, "Body did not match expected output")
	})

	t.Run("Linux arguments are passed through to end script", func(t *testing.T) {
		if runtime.GOOS == "linux" {
			configuration.Settings.Security.ApprovedExecutablesOnly.IsTrue = false
			configuration.Settings.Security.SignedStdInOnly.IsTrue = false
			configuration.Settings.Security.AllowScriptArguments.IsTrue = false

			testRequest := map[string]interface{}{
				"path":  `sh`,
				"args":  []string{"-s"},
				"stdin": `#!/bin/bash\necho "First line.";\necho "Second line.";\necho "Third lime.";\n`,
				"stdinsignature": `untrusted comment: signature from minisign secret key
RWQ3ly9IPenQ6Yds9dwf0ZbgW3nOe6pwhgdaFPeoXSO8eNNInMRE5UEe+lsGuWG016SeNAbWKtuZVOV5QxBcuTNukoMmB8+z7A0=
trusted comment: timestamp:1631633837	file:echo.sh
XSx0EYiti1RA1IEQd7HCLyF0cEEKj5xXSmiKV9BnPmrRHKcS5Et35Xhpynlu0t1RLlZUQDueRqAwmyaunxjqAw==				
`,
			}

			expectedOutput := `First line.\nSecond line.\nThird lime.\n`

			output := RunTestRequest(t, http.MethodPost, "/v1/runscriptstdin", JsonSerialize(testRequest))

			assert := assert.New(t)
			assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
			assert.Equal(expectedOutput, output.ResponseBody)
		}
	})
}
