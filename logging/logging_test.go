package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLogger(t *testing.T) {
	t.Cleanup(func() { SetLogger(nil) })
	var output bytes.Buffer
	SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))
	Info("started", "strategy", "example")
	Error("failed", "error", "test error")
	for _, want := range []string{`"msg":"started"`, `"strategy":"example"`, `"msg":"failed"`, `"error":"test error"`} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("log output %q does not contain %q", output.String(), want)
		}
	}

	var defaultOutput bytes.Buffer
	standard := logrus.StandardLogger()
	oldOutput := standard.Out
	standard.SetOutput(&defaultOutput)
	t.Cleanup(func() { standard.SetOutput(oldOutput) })
	SetLogger(nil)
	Error("default", "error", "test error")
	if !strings.Contains(defaultOutput.String(), "default") || !strings.Contains(defaultOutput.String(), "test error") {
		t.Errorf("default logger output: %q", defaultOutput.String())
	}
}

type countingLogger struct {
	count atomic.Int64
}

func (l *countingLogger) Info(string, ...any)  { l.count.Add(1) }
func (l *countingLogger) Error(string, ...any) { l.count.Add(1) }

func TestConcurrentLoggingAndReplacement(t *testing.T) {
	t.Cleanup(func() { SetLogger(nil) })
	l := &countingLogger{}
	SetLogger(l)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				SetLogger(l)
				Info("event")
				Error("event")
			}
		}()
	}
	wg.Wait()
	if got := l.count.Load(); got != 2000 {
		t.Fatalf("got %d events, want 2000", got)
	}
}
