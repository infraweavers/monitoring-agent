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
		MaxBackups: configuration.Settings.LogArchiveFilesToRetain,
		MaxSize:    configuration.Settings.LogRotationThresholdInMegaBytes,
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

// LogCriticalf records a CRITICAL log message formatted according to a format specifier if log configuration.Settings.LogLevel >= CRITICAL
func LogCriticalf(message string, v ...interface{}) {
	writef(CRITICAL, message, v...)
}

// LogCritical records a CRITCAL log message if log configuration.Settings.LogLevel >= CRITICAL
func LogCritical(message string) {
	write(CRITICAL, message)
}

// LogErrorf records an ERROR log message formatted according to a format specifier if log configuration.Settings.LogLevel >= ERROR
func LogErrorf(message string, v ...interface{}) {
	writef(ERROR, message, v...)
}

// LogError records a ERROR log message if log configuration.Settings.LogLevel >= ERROR
func LogError(message string) {
	write(ERROR, message)
}

// LogWarningf records a WARNING log message formatted according to a format specifier if log configuration.Settings.LogLevel >= WARNING
func LogWarningf(message string, v ...interface{}) {
	writef(WARNING, message, v...)
}

// LogWarning records a WARNING log message if log configuration.Settings.LogLevel >= WARNING
func LogWarning(message string) {
	write(WARNING, message)
}

// LogNoticef records a NOTICE log message formatted according to a format specifier if log configuration.Settings.LogLevel >= NOTICE
func LogNoticef(message string, v ...interface{}) {
	writef(NOTICE, message, v...)
}

// LogNotice records a NOTICE log message if log configuration.Settings.LogLevel >= NOTICE
func LogNotice(message string) {
	write(NOTICE, message)
}

// LogInfof records a INFO log message formatted according to a format specifier if log configuration.Settings.LogLevel >= INFO
func LogInfof(message string, v ...interface{}) {
	writef(INFO, message, v...)
}

// LogInfo records a INFO log message if log configuration.Settings.LogLevel >= INFO
func LogInfo(message string) {
	write(INFO, message)
}

// LogDebugf records a DEBUG log message formatted according to a format specifier if log configuration.Settings.LogLevel >= DEBUG
func LogDebugf(message string, v ...interface{}) {
	writef(DEBUG, message, v...)
}

// LogDebug records a DEBUG log message if log configuration.Settings.LogLevel >= DEBUG
func LogDebug(message string) {
	write(DEBUG, message)
}
