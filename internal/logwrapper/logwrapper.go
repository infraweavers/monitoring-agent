package logwrapper

import (
	"errors"
	"fmt"
	"monitoringagent/internal/configuration"
	"os"
	"strings"

	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Log wrapped instance
//var Log = logging.MustGetLogger("default")

type logLevel int

const (
	CRITICAL logLevel = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

func (p logLevel) String() string {
	return levelNames[p]
}

func parseLogLevel(level string) (logLevel, error) {
	for i, name := range levelNames {
		if strings.EqualFold(name, level) {
			return logLevel(i), nil
		}
	}
	return ERROR, errors.New("invalid log level")
}

var level logLevel
var Log *log.Logger
var RunningInteractively = false

// Initialise and configure the logging
func Initialise(runningInteractively bool) {

	logFile, err := os.OpenFile(configuration.Settings.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	level, err = parseLogLevel(configuration.Settings.LogLevel)
	if err != nil {
		panic(err)
	}

	Log = log.New(logFile, "ma ", log.LstdFlags)

	Log.SetOutput(&lumberjack.Logger{
		Filename:   logFile.Name(),
		MaxBackups: 10,
		MaxAge:     1,
	})

	RunningInteractively = runningInteractively

	LogInfo("Logging Initialised")
	LogInfof("Logging with LogFilePath: '%s'", configuration.Settings.LogFilePath)
	LogInfof("Running Interactively: %t", runningInteractively)
}

func writef(lvl logLevel, message string, v ...interface{}) {
	if lvl > level {
		return
	}
	format := fmt.Sprintf("%s : %s", logLevel.String(level), message)
	Log.Printf(format, v...)

	if RunningInteractively {
		log.Printf(message, v...)
	}
}

func write(lvl logLevel, message string) {
	if lvl > level {
		return
	}
	message = logLevel.String(level) + ": " + message
	Log.Print(message)

	if RunningInteractively {
		log.Println(message)
	}
}

func LogCriticalf(message string, v ...interface{}) {
	writef(CRITICAL, message, v...)
}

func LogCritical(message string) {
	write(CRITICAL, message)
}

func LogErrorf(message string, v ...interface{}) {
	writef(ERROR, message, v...)
}

func LogError(message string) {
	write(ERROR, message)
}

func LogWarningf(message string, v ...interface{}) {
	writef(WARNING, message, v...)
}

func LogWarning(message string) {
	write(WARNING, message)
}

func LogNoticef(message string, v ...interface{}) {
	writef(NOTICE, message, v...)
}

func LogNotice(message string) {
	write(NOTICE, message)
}

func LogInfof(message string, v ...interface{}) {
	writef(INFO, message, v...)
}

func LogInfo(message string) {
	write(INFO, message)
}

func LogDebugf(message string, v ...interface{}) {
	writef(DEBUG, message, v...)
}

func LogDebug(message string) {
	write(DEBUG, message)
}
