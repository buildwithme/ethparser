package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log.New(os.Stdout, "[CLI] ", log.LstdFlags),
	}
}
