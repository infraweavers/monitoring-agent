package logwrapper

import (
	"mama/internal/configuration"
	"os"

	"github.com/op/go-logging"
)

// Log wrapped instance
var Log = logging.MustGetLogger("default")

// Initialise Configure the logging
func Initialise() {

	logFile, logError := os.OpenFile(configuration.Settings.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	logging.LogLevel(configuration.Settings.LogLevel)

	if logError != nil {
		Log.Fatal(logError)
	}

	logging.SetBackend(logging.NewLogBackend(logFile, "", 0))

	Log.Info("Logging Initialised")
}
