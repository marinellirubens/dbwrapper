package logger

import (
	"fmt"
	"log"
	"os"
)

// log level
const (
	NONE    = 0
	DEBUG   = 10
	INFO    = 20
	WARNING = 30
	ERROR   = 40
	FATAL   = 50
)

// log type
const (
	//writes the lods on the terminal
	STREAM_WRITER = 777

	//writes the logs on a file
	FILE_WRITER = 999
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

// TODO: Implement log rotation
type Logger struct {
	logger   *log.Logger
	logLevel int
}

// generic log method
func (l *Logger) log(message string, logLevel int) {
	fmt.Println(message, logLevel)
}

// logs message using level Info
func (l *Logger) Info(message string) {
	if l.logLevel >= INFO {
		l.log(message, INFO)
	}
}

// logs message using level Debug
func (l *Logger) Debug(message string) {
	if l.logLevel >= DEBUG {
		l.log(message, DEBUG)
	}
}

// logs message using level Warning
func (l *Logger) Warning(message string) {
	if l.logLevel >= WARNING {
		l.log(message, WARNING)
	}
}

// logs message using level Error
func (l *Logger) Error(message string) {
	if l.logLevel >= ERROR {
		l.log(message, ERROR)
	}
}

// logs message using level Fatal
func (l *Logger) Fatal(message string) {
	if l.logLevel >= FATAL {
		l.log(message, FATAL)
	}
}

// create logger for a more easy logging process
// the logging level shoud be based on the constants of the package
// also the logType should be base on the constants of the packag
//
// example:
//
//	logger, err := lCreateLogger("server.log", logger.DEBUG, logger.STREAM_WRITER)
func CreateLogger(logFile string, logLevel int, logType int) (*Logger, error) {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	var output *os.File
	var err error

	if logType == FILE_WRITER {
		file, err := os.Create(logFile)
		if err != nil {
			defer file.Close()
			return nil, err
		}
		output = file
	} else if logType == STREAM_WRITER {
		output = os.Stdout
	} else {
		panic("Log type invalid")
	}
	var internalLogger *log.Logger

	internalLogger = log.New(output, "INFO: ", flags)

	logger := &Logger{
		logger:   internalLogger,
		logLevel: logLevel}
	return logger, err
}
