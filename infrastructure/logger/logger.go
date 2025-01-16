package logger

import (
	"log"
	"os"
)

type Logger struct {
}

type LoggerInterface interface {
	LogError(err error)
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogError(err error) {
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	errorLog.Printf("error: %s", err.Error())
}
