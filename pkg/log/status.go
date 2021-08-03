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
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

const (
	successStatusFormat = " \x1b[32m✓\x1b[0m %s"
	failureStatusFormat = " \x1b[31m✗\x1b[0m %s"
)

var (
	defaultCharSet = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	defaultDelay   = 100 * time.Millisecond
)

type Status struct {
	spinner *spinner.Spinner
	message string
	logger  *Logger

	successStatusFormat string
	failureStatusFormat string
}

func NewStatus(logger *Logger) *Status {
	return &Status{
		logger:              logger,
		spinner:             spinner.New(defaultCharSet, defaultDelay, spinner.WithWriter(os.Stdout)),
		successStatusFormat: successStatusFormat,
		failureStatusFormat: failureStatusFormat,
	}
}

func (s *Status) Start(msg string) {
	s.Stop()
	s.message = msg
	if s.spinner != nil {
		s.spinner.Suffix = fmt.Sprintf(" %s ", s.message)
		s.spinner.Start()
	}
}

func (s *Status) Success(msg ...string) {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		s.spinner.Stop()
	}
	if msg != nil {
		s.logger.Infof(s.successStatusFormat, strings.Join(msg, " "))
		return
	}
	s.logger.Infof(s.successStatusFormat, s.message)
	s.message = ""
}

func (s *Status) Error(msg ...string) {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		s.spinner.Stop()
	}
	if msg != nil {
		s.logger.Infof(s.failureStatusFormat, strings.Join(msg, " "))
		return
	}
	s.logger.Infof(s.failureStatusFormat, s.message)
	s.message = ""
}

func (s *Status) Stop() {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		fmt.Fprint(s.logger.Writer, "\r")
		s.spinner.Stop()
	}
}
