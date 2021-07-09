package log

import (
	"testing"
	"time"
)

func TestSuccessMessage(t *testing.T) {
	s := NewStatus()
	l := NewLogger(true)
	s.Start("Installing app1")
	time.Sleep(2 * time.Second)
	l.Info("Waiting for pods to be ready!!")
	l.Debug("Waiting for pods to be ready!!")
	time.Sleep(3 * time.Second)
	s.End(true)
	time.Sleep(time.Second)
	s.Start("Installing app2")
	time.Sleep(5 * time.Second)
	s.End(false)
}
