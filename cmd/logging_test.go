package bot

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBotLogger(t *testing.T) {
	oldConfigFile := GlobalFlags.ConfigFile
	t.Cleanup(func() {
		GlobalFlags.ConfigFile = oldConfigFile
		SetLogger(nil)
	})
	var output bytes.Buffer
	SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))

	GlobalFlags.ConfigFile = filepath.Join(t.TempDir(), "missing.yaml")
	executeStartCommand(nil, nil)
	if !strings.Contains(output.String(), `"msg":"Cannot load bot configuration"`) {
		t.Fatalf("configuration error not logged: %s", output.String())
	}

	output.Reset()
	GlobalFlags.ConfigFile = filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(GlobalFlags.ConfigFile, []byte("simulation_mode: true\n"), 0600); err != nil {
		t.Fatal(err)
	}
	executeStartCommand(nil, nil)
	for _, want := range []string{`"msg":"Starting bot"`, `"msg":"Bot stopped"`} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("log output %q does not contain %q", output.String(), want)
		}
	}
}
