package logwrapper

import (
	"mama/internal/configuration"
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

	if runningInteractively {
		logging.SetBackend(logging.MultiLogger(logging.NewLogBackend(logFile, "", 0), logging.NewLogBackend(os.Stderr, "", 0)))
	} else {
		logging.SetBackend(logging.NewLogBackend(logFile, "", 0))
	}

	Log.Info("Logging Initialised")
	Log.Infof("Logging with LogFilePath: '%s'", configuration.Settings.LogFilePath)
	Log.Infof("Running Interactively: %t", runningInteractively)
}
