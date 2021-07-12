package log

import (
	"testing"
	"time"
)

func TestSuccessMessage(t *testing.T) {
	l := NewLogger(true)
	s := NewStatus(l)
	l.Info("Initializing")
	s.Start("Setting up deps")
	time.Sleep(2 * time.Second)
	l.Debug("Waiting for pods to be ready!!")
	time.Sleep(time.Second)
	s.Stop()
	s.Start("Installing app2")
	time.Sleep(2 * time.Second)
	s.Success()
	s.Start("Installing app3")
	time.Sleep(2 * time.Second)
	l.Error("Timed out while waiting for pods to be ready!!")
	s.Error("App3 Installation failed!")
}
