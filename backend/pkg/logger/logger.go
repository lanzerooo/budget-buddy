package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

var logger *Logger

func Init() {
	logger = &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(v ...interface{}) {
	logger.infoLogger.Println(v...)
}

func Warn(v ...interface{}) {
	logger.warnLogger.Println(v...)
}

func Error(v ...interface{}) {
	logger.errorLogger.Println(v...)
}

func Fatal(v ...interface{}) {
	logger.errorLogger.Println(v...)
	os.Exit(1)
}

func Printf(format string, v ...interface{}) {
	logger.infoLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.errorLogger.Printf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	logger.warnLogger.Printf(format, v...)
}
