package strategies

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/logging"
)

func TestStrategyLogging(t *testing.T) {
	var output bytes.Buffer
	logging.SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() { logging.SetLogger(nil) })
	wantErr := errors.New("test failure")
	callbackCalled := false
	strategy := IntervalStrategy{Model: StrategyModel{
		Name: "test",
		OnUpdate: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			return wantErr
		},
		OnError: func(err error) { callbackCalled = err == wantErr },
	}}
	(&Tactic{Strategy: strategy}).Execute(nil)
	if !callbackCalled {
		t.Fatal("strategy error callback was not called")
	}
	for _, want := range []string{`"msg":"Starting strategy"`, `"msg":"Strategy update failed"`, `"strategy":"test"`, `"error":"test failure"`, `"msg":"Strategy stopped"`} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("log output %q does not contain %q", output.String(), want)
		}
	}
}
