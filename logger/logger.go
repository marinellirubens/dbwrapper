package logger

import (
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
	STREAM_WRITER = 777
	FILE_WRITER   = 999
)

// TODO: Implement log rotation
type Logger struct {
	infoLogger    *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
}

func (l *Logger) Info(message string) {
	if l.infoLogger != nil {
		l.infoLogger.Print(message)
	}
}

func (l *Logger) Debug(message string) {
	if l.debugLogger != nil {
		l.debugLogger.Print(message)
	}
}

func (l *Logger) Warning(message string) {
	if l.warningLogger != nil {
		l.warningLogger.Print(message)
	}
}

func (l *Logger) Error(message string) {
	if l.errorLogger != nil {
		l.errorLogger.Print(message)
	}
}

func (l *Logger) Fatal(message string) {
	if l.fatalLogger != nil {
		l.fatalLogger.Print(message)
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
	var info *log.Logger
	var debug *log.Logger
	var warning *log.Logger
	var error_ *log.Logger
	var fatal *log.Logger

	if logLevel <= INFO {
		info = log.New(output, "INFO: ", flags)
	}
	if logLevel <= DEBUG {
		debug = log.New(output, "DEBUG: ", flags)
	}
	if logLevel <= WARNING {
		warning = log.New(output, "WARNING: ", flags)
	}
	if logLevel <= ERROR {
		error_ = log.New(output, "ERROR: ", flags)
	}
	if logLevel <= FATAL {
		fatal = log.New(output, "FATAL: ", flags)
	}

	logger := &Logger{
		infoLogger:    info,
		debugLogger:   debug,
		warningLogger: warning,
		errorLogger:   error_,
		fatalLogger:   fatal}
	return logger, err
}
