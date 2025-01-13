package logger

import (
	"log"
	"os"
)

type LoggerInterface interface {
	LogError(statusCode int, err error)
}

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) LogError(statusCode int, err error) {
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	errorLog.Printf("status code: %d, error: %s", statusCode, err.Error())
}
