package logging

import (
	"mama/internal/configuration"
	"os"

	"github.com/sirupsen/logrus"
)

// Log Logrus instance
var Log = logrus.New()

// Initialise Configure the logging
func Initialise() {

	logFile, logError := os.OpenFile(configuration.Settings.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	level, levelError := logrus.ParseLevel(configuration.Settings.LogLevel)

	if logError != nil {
		Log.Fatal(logError)
	}

	if levelError != nil {
		Log.Fatal(levelError)
	}

	Log.SetLevel(level)

	Log.Out = logFile
	Log.Info("Logging Initialised")
}
