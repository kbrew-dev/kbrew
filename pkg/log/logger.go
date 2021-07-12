package log

import (
	"fmt"
	"os"
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
	fmt.Fprint(os.Stdout, "\r")
	fmt.Println(message...)
}

func (l *Logger) Infof(format string, message ...interface{}) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf(format, message...)
}

func (l *Logger) Debug(message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("DEBUG: ")
	fmt.Println(message...)
}

func (l *Logger) Debugf(format string, message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("DEBUG: "+format, message...)
}

func (l *Logger) Error(message ...interface{}) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("\x1b[31mERROR:\x1b[0m ")
	fmt.Println(message...)
}

func (l *Logger) Errorf(format string, message ...interface{}) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("\x1b[31mERROR:\x1b[0m "+format, message...)
}

func (l *Logger) Warn(message ...interface{}) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("\x1b[33mWARN:\x1b[0m ")
	fmt.Println(message...)
}

func (l *Logger) Warnf(format string, message ...interface{}) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("\x1b[33mWARN:\x1b[0m "+format, message...)
}

func (l *Logger) InfoMap(key, value string) {
	fmt.Fprint(os.Stdout, "\r")
	fmt.Printf("\x1b[32m%s:\x1b[0m %s\n", key, value)
}
