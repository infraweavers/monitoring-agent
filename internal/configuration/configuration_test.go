package configuration

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var paths Paths

func test(jsonbytes []byte) (settings Config, err error) {
	workingDirectory, _ := os.Getwd()
	paths = InitialisePaths(workingDirectory)

	settings = Config{}
	err = json.Unmarshal(jsonbytes, &settings)
	settings.Paths.Reset(paths)

	return
}
func TestConfigJsonImport(t *testing.T) {

	settings, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.Empty(t, err, "unmarshalling json should not result in an error")
	assert.NotEmpty(t, settings, "settings struct is populated with no missing mandatory values")
	assert.NotEmpty(t, settings.Paths, "settings.Paths struct is populated")
	assert.NotEmpty(t, settings.Authentication.Password, "settings.Authentication.Password is populated")
	assert.Equal(t, true, settings.Logging.LogHTTPRequests.HasValue, "NullBool type (settings.Logging.LogHTTPRequests) is not null")
	assert.Equal(t, ConfigurationDirectory+"\\"+"output.log", settings.Logging.LogFilePath, "settings.Logging.LogFilePath includes ConfigurationDirectory %s", ConfigurationDirectory)
	assert.Equal(t, ConfigurationDirectory+"\\"+"PathToClientCertificateCAFile", settings.Security.ClientCertificateCAFile.Path, "settings.Security.ClientCertificateCAFile is set and path includes ConfigurationDirectory %s", ConfigurationDirectory)

	expectedDuration, _ := time.ParseDuration("15s")
	assert.Equal(t, expectedDuration, settings.Server.DefaultScriptTimeout.Duration, "settings.Server.DefaultScriptTimeout matches supplied JSON value")

}

func TestLogFilesToRetainSetToZeroDoesNotCauseError(t *testing.T) {
	settings, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 0,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.Empty(t, err, "Setting LogFilesToRetain to 0 should not cause an error on account of zero value being treated as null")
	assert.Equal(t, true, settings.Logging.LogArchiveFilesToRetain.HasValue, "settings.Logging.LogArchiveFilesToRetain.HasValue should be true even if LogFilesToRetain was set to 0")
	assert.Equal(t, 0, settings.Logging.LogArchiveFilesToRetain.Value, "Setting LogFilesToRetain to 0 should result in settings.Logging.LogArchiveFilesToRetain also being 0")
}

func TestConfigJsonImportFailsWhenAuthorisationSectionIsMissing(t *testing.T) {
	_, err := test([]byte(`
	{
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "Authentication not set when tagged with 'mandatory:\"true\"'", "Error returned when configuration authorisation section is missing")
}

func TestConfigJsonImportFailsWhenAuthorisationUsernameIsMissing(t *testing.T) {

	_, err := test([]byte(`
	{
		"Authentication": {
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "Username not set when tagged with 'mandatory:\"true\"'", "Error returned when authorisation username is omitted")
}

func TestConfigJsonImportFailsWhenLogArchiveFilesToRetainIsMissing(t *testing.T) {

	_, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "LogArchiveFilesToRetain not set when tagged with 'mandatory:\"true\"'", "Error returned when Logging LogArchiveFilesToRetain is omitted")
}

func TestConfigJsonImportFailsWhenLogHTTPRequestsIsMissing(t *testing.T) {

	_, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "LogHTTPRequests not set when tagged with 'mandatory:\"true\"'", "Error returned when Logging LogHTTPRequests is omitted")
}

func TestConfigJsonImportFailsWhenHTTPRequestTimeoutIsMissing(t *testing.T) {

	_, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ClientCertificateCAFile": "PathToClientCertificateCAFile",
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "HTTPRequestTimeout not set when tagged with 'mandatory:\"true\"'", "Error returned when Server HTTPRequestTimeout is omitted")
}

func TestConfigJsonImportFailsWhenClientCertificateCAFileIsMissing(t *testing.T) {

	_, err := test([]byte(`
	{
		"Authentication": {
			"Username": "test",
			"Password": "secret"
		},
		"Logging": {
			"LogFilePath": "output.log",
			"LogLevel": "INFO",
			"LogArchiveFilesToRetain": 10,
			"LogRotationThresholdInMegaBytes": 100,
			"LogHTTPRequests": false,
			"LogHTTPResponses": false
		},
		"Server": {
			"BindAddress": "0.0.0.0:9000",
			"HTTPRequestTimeout": "300s",
			"DefaultScriptTimeout": "15s",        
			"LoadPprof": false
		},
		"Security": {
			"DisableHTTPs": false,
			"SignedStdInOnly": false,
			"PublicKey": "RWTV8L06+shYI7Xw1H+NBGmsUYlbEkbrdYxr4c0ImLCAr8NGx75VhxGQ",
			"AllowedAddresses": ["::1/128","127.0.0.0/8","0.0.0.0/0"],
			"UseClientCertificates": false,
			"ApprovedPathArgumentsOnly": false,
			"ApprovedPathArguments": {
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe": [
					["-command", "-"],
					["-command","start-sleep 1"],
					["-command","Write-Host 'Hello, World'"],
					["-command","Write-Host \"Hello, World\""],
					["-command"]
				],
				"sh": [
					["-s"],
					["-c"]
				]
			}
		}
	}
	`))

	assert.EqualError(t, err, "ClientCertificateCAFile not set when tagged with 'mandatory:\"true\"'", "Error returned when Security ClientCertificateCAFile is omitted")
}
