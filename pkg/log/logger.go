// Copyright 2021 The kbrew Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"io"
	"os"
)

const (
	errorPrefix = "\x1b[31mERROR:\x1b[0m "
	warnPrefix  = "\x1b[33mWARN:\x1b[0m "
	debugPrefix = "DEBUG: "

	infoMapKeyFormat = "\x1b[32m%s:\x1b[0m "
)

type Logger struct {
	DebugLevel bool
	Writer     io.Writer
}

func NewLogger(debug bool) *Logger {
	return &Logger{
		DebugLevel: debug,
		Writer:     os.Stdout,
	}
}

func (l *Logger) SetWriter(writer io.Writer) {
	l.Writer = writer
}

func (l *Logger) Info(message ...interface{}) {
	l.print("", message...)
}

func (l *Logger) Infof(format string, message ...interface{}) {
	l.print("", fmt.Sprintf(format, message...))
}

func (l *Logger) Debug(message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	l.print(debugPrefix, message...)
}

func (l *Logger) Debugf(format string, message ...interface{}) {
	if !l.DebugLevel {
		return
	}
	l.print(debugPrefix, fmt.Sprintf(format, message...))
}

func (l *Logger) Error(message ...interface{}) {
	l.print(errorPrefix, message...)
}

func (l *Logger) Errorf(format string, message ...interface{}) {
	l.print(errorPrefix, fmt.Sprintf(format, message...))
}

func (l *Logger) Warn(message ...interface{}) {
	l.print(warnPrefix, message...)
}

func (l *Logger) Warnf(format string, message ...interface{}) {
	l.print(warnPrefix, fmt.Sprintf(format, message...))
}

func (l *Logger) InfoMap(key, value string) {
	l.print(fmt.Sprintf(infoMapKeyFormat, key), value)
}

func (l *Logger) print(prefix string, message ...interface{}) {
	fmt.Fprint(l.Writer, "\r")
	fmt.Fprint(l.Writer, prefix)
	fmt.Fprintln(l.Writer, message...)
}
