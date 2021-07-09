package log

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	//"k8s.io/klog/v2"
)

const (
	successStatusFormat = " \x1b[32m✓\x1b[0m %s\n"
	failureStatusFormat = " \x1b[31m✗\x1b[0m %s\n"
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
	s.End(true)
	s.message = msg
	if s.spinner != nil {
		s.spinner.Suffix = fmt.Sprintf(" %s ", s.message)
		s.spinner.Start()
	}
}

func (s *Status) End(success bool) {
	if s.message == "" {
		return
	}
	if s.spinner != nil {
		s.spinner.Stop()
		fmt.Fprint(s.spinner.Writer, "\r")
	}
	if success {
		s.logger.Infof(s.successStatusFormat, s.message)
	} else {
		s.logger.Infof(s.failureStatusFormat, s.message)
	}
	s.message = ""
}
