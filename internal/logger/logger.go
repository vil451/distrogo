package logger

import "log"

const (
	LogLevelDebug   = iota
	LogLevelInfo    = iota
	LogLevelWarning = iota
	LogLevelError   = iota
)

var logLevel int

func SetLogLevel(level int) {
	logLevel = level
}

func Debug(v ...interface{}) {
	if logLevel <= LogLevelDebug {
		log.Print(v...)
	}
}

func Info(v ...interface{}) {
	if logLevel <= LogLevelInfo {
		log.Print(v...)
	}
}

func Warning(v ...interface{}) {
	if logLevel <= LogLevelWarning {
		log.Print(v...)
	}
}

func Error(v ...interface{}) {
	if logLevel <= LogLevelError {
		log.Print(v...)
	}
}
