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
		s.spinner.Stop()
	}
}
