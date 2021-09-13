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

// LogLevel constant, CRITICAL = 0, DEBUG = 5
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

// Log exposes the log packages logger for functions such as FatalF and PanicF
var Log *log.Logger

var runningInteractively = false
var linebreak string

// Initialise and configure the logging
func Initialise(isRunInteractively bool, newline string) {

	logFile, err := os.OpenFile(configuration.Settings.Logging.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	level, err = parseLogLevel(configuration.Settings.Logging.LogLevel)
	if err != nil {
		panic(err)
	}

	linebreak = newline
	Log = log.New(logFile, "ma ", log.LstdFlags)

	Log.SetOutput(&lumberjack.Logger{
		Filename:   logFile.Name(),
		MaxBackups: configuration.Settings.Logging.LogArchiveFilesToRetain.Value,
		MaxSize:    configuration.Settings.Logging.LogRotationThresholdInMegaBytes,
	})

	runningInteractively = isRunInteractively

	LogInfo("Logging Initialised")
	LogInfof("Logging with LogFilePath: '%s'", configuration.Settings.Logging.LogFilePath)
	LogInfof("Running Interactively: %t", runningInteractively)
}

func writef(lvl logLevel, message string, v ...interface{}) {
	if lvl > level {
		return
	}
	format := fmt.Sprintf("%s : %s%s", logLevel.String(lvl), message, linebreak)
	Log.Printf(format, v...)

	if runningInteractively {
		log.Println(fmt.Sprintf(format, v...))
	}
}

func write(lvl logLevel, message string) {
	if lvl > level {
		return
	}
	message = fmt.Sprintf("%s : %s%s", logLevel.String(lvl), message, linebreak)
	Log.Print(message)

	if runningInteractively {
		log.Println(message)
	}
}

// LogHTTPRequest records HTTP Request information in the log
func LogHTTPRequest(remoteAddr string, host string, method string, url string, header map[string][]string, proto string, contentLength int64, body string) {
	Log.Printf("HTTP Request: %s [%s] %s %#v %s %d %s", remoteAddr, method, url, header, proto, contentLength, body)

	if runningInteractively {
		log.Printf("HTTP Request: %s [%s] %s %#v %s %d %s", remoteAddr, method, url, header, proto, contentLength, body)
	}
}

// LogHTTPResponse records HTTP response information in the log
func LogHTTPResponse(status string, header map[string][]string, proto string, body string) {
	Log.Printf("HTTP Response: %s %#v %s %s", status, header, proto, body)

	if runningInteractively {
		log.Printf("HTTP Response: %s %#v %s %s", status, header, proto, body)
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
