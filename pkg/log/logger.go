package log

import (
	"fmt"
)

type Logger struct {
	DebugLevel bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{
		DebugLevel: debug,
	}
}

func (l *Logger) Info(message ...interface{}) {
	fmt.Println(message...)
}

func (l *Logger) Infof(format string, message ...interface{}) {
	fmt.Printf(format, message...)
}

func (l *Logger) Debug(message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	fmt.Printf("DEBUG: ")
	fmt.Println(message...)
}

func (l *Logger) Debugf(format string, message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	fmt.Printf("DEBUG: "+format, message...)
}
