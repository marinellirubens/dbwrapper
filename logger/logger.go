package logger

import (
	"fmt"
	"log"
	"os"
)

// log level
const (
	NONE    = uint8(0)
	DEBUG   = uint8(10)
	INFO    = uint8(20)
	WARNING = uint8(30)
	ERROR   = uint8(40)
	FATAL   = uint8(50)
)

// log type
const (
	// writes the logs on the terminal
	STREAM_WRITER = uint16(777)

	// writes the logs on a file
	FILE_WRITER = uint16(999)
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

var (
	levelTxt map[uint8]string = map[uint8]string{
		uint8(0):  "NONE",
		uint8(10): "DEBUG",
		uint8(20): "INFO",
		uint8(30): "WARNING",
		uint8(40): "ERROR",
		uint8(50): "FATAL",
	}

	levelTxtWithColor map[uint8]string = map[uint8]string{
		uint8(0):  "NONE",
		uint8(10): green + "DEBUG" + reset,
		uint8(20): blue + "INFO" + reset,
		uint8(30): yellow + "WARNING" + reset,
		uint8(40): red + "ERROR" + reset,
		uint8(50): red + "FATAL" + reset,
	}
)

type Handler struct {
	logger   *log.Logger
	logLevel uint8
	logType  uint16
}

// TODO: Implement log rotation
type Logger struct {
	handlers []*Handler
}

// converts log level to text with color
func levelToText(logLevel uint8, withColor bool) string {
	var ret string
	if withColor {
		ret = fmt.Sprint(levelTxtWithColor[logLevel])
	} else {
		ret = fmt.Sprint(levelTxt[logLevel])
	}
	return ret
}

// generic log method
func (l *Handler) log(message string, logLevel uint8) {
	var withColor bool
	if l.logType == STREAM_WRITER {
		withColor = true
	}
	level := levelToText(logLevel, withColor)
	text_to_print := fmt.Sprintf("[DBWRAPPER][%s] - %s%s", level, message, reset)
	l.logger.Println(text_to_print)
}

// logs message using level Info
func (l *Logger) Info(message string) {
	for _, handler := range l.handlers {
		if handler.logLevel <= INFO {
			handler.log(message, INFO)
		}
	}
}

// logs message using level Debug
func (l *Logger) Debug(message string) {
	for _, handler := range l.handlers {
		if handler.logLevel <= DEBUG {
			handler.log(message, DEBUG)
		}
	}
}

// logs message using level Warning
func (l *Logger) Warning(message string) {
	for _, handler := range l.handlers {
		if handler.logLevel <= WARNING {
			handler.log(message, WARNING)
		}
	}
}

// logs message using level Error
func (l *Logger) Error(message string) {
	for _, handler := range l.handlers {
		if handler.logLevel <= ERROR {
			handler.log(message, ERROR)
		}
	}
}

// logs message using level Fatal
func (l *Logger) Fatal(message string) {
	for _, handler := range l.handlers {
		if handler.logLevel <= FATAL {
			handler.log(message, FATAL)
		}
	}
}

// create logger for a more easy logging process
// the logging level shoud be based on the constants of the package
// also the logType should be base on the constants of the packag
//
// example:
//
//	logger, err := CreateLogger("server.log", logger.DEBUG, logger.STREAM_WRITER)
func CreateLogger(logFile string, logLevel uint8, logTypes []uint16) (*Logger, error) {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	var output *os.File
	var err error
	var handlers []*Handler
	var internalLogger *log.Logger

	for _, logType := range logTypes {
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
		internalLogger = log.New(output, "", flags)
		handler := &Handler{
			logger:   internalLogger,
			logLevel: logLevel,
			logType:  logType,
		}
		handlers = append(handlers, handler)
	}

	logger := &Logger{
		handlers: handlers,
	}
	return logger, err
}
