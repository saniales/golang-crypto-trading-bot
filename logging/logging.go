// Package logging provides the bot's shared operational logger.
package logging

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger is implemented by *slog.Logger and can be adapted to other logging systems.
// Implementations must be safe for concurrent use.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

var (
	mu     sync.RWMutex
	logger Logger = logrusLogger{}
)

// SetLogger replaces the operational logger. Passing nil restores the default Logrus logger.
func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	if l == nil {
		l = logrusLogger{}
	}
	logger = l
}

// Info records an informational event.
func Info(msg string, args ...any) {
	mu.RLock()
	l := logger
	mu.RUnlock()
	l.Info(msg, args...)
}

// Error records an error event.
func Error(msg string, args ...any) {
	mu.RLock()
	l := logger
	mu.RUnlock()
	l.Error(msg, args...)
}

type logrusLogger struct{}

func (logrusLogger) Info(msg string, args ...any) {
	logrus.WithFields(fields(args)).Info(msg)
}

func (logrusLogger) Error(msg string, args ...any) {
	logrus.WithFields(fields(args)).Error(msg)
}

func fields(args []any) logrus.Fields {
	f := make(logrus.Fields, (len(args)+1)/2)
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprint(args[i])
		if i+1 < len(args) {
			f[key] = args[i+1]
		} else {
			f[key] = nil
		}
	}
	return f
}
