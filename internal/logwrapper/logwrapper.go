package logwrapper

import (
	"monitoringagent/internal/configuration"
	"os"

	"github.com/op/go-logging"
)

// Log wrapped instance
var Log = logging.MustGetLogger("default")

// Initialise Configure the logging
func Initialise(runningInteractively bool, configurationDirectory string) {

	logFile, logError := os.OpenFile(configuration.Settings.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	logging.LogLevel(configuration.Settings.LogLevel)

	if logError != nil {
		Log.Fatal(logError)
	}

	logFileBackend := logging.NewLogBackend(logFile, "", 0)
	backendFormatter := logging.NewBackendFormatter(logFileBackend, logging.MustStringFormatter(`%{time}: %{message}`))
	if runningInteractively {
		logging.SetBackend(logging.MultiLogger(backendFormatter, logging.NewLogBackend(os.Stderr, "", 0)))
	} else {
		logging.SetBackend(backendFormatter)
	}

	Log.Info("Logging Initialised")
	Log.Infof("Logging with LogFilePath: '%s'", configuration.Settings.LogFilePath)
	Log.Infof("Running Interactively: %t", runningInteractively)
}
