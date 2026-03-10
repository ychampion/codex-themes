package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/cli"
)

func TestSetupCommandAppliesThemeAndWindowsTerminal(t *testing.T) {
	t.Parallel()

	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Abs() error = %v", err)
	}
	codexHome := filepath.Join(t.TempDir(), ".codex")
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(settingsPath, []byte(`{"profiles":{"defaults":{}},"schemes":[]}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := cli.NewRootCommand(cli.Options{ThemeDirs: []string{filepath.Join(rootDir, "themes")}})
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"setup", "--theme", "amber-nocturne", "--codex-home", codexHome, "--windows-terminal-settings", settingsPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\nstderr:\n%s", err, stderr.String())
	}

	configBytes, err := os.ReadFile(filepath.Join(codexHome, "config.toml"))
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	configText := string(configBytes)
	for _, fragment := range []string{"[codex_themes]", "[tui]", "theme = 'codex-themes-amber-nocturne'", "model-with-reasoning", "current-dir", "git-branch"} {
		if !strings.Contains(configText, fragment) {
			t.Fatalf("config missing fragment %q\nconfig:\n%s", fragment, configText)
		}
	}

	if _, err := os.Stat(filepath.Join(codexHome, "themes", "codex-themes-amber-nocturne.tmTheme")); err != nil {
		t.Fatalf("expected codex theme file: %v", err)
	}

	if !strings.Contains(configText, "[codex_themes]") {
		t.Fatalf("config missing managed section\nconfig:\n%s", string(configBytes))
	}

	settingsBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(settings) error = %v", err)
	}
	if !strings.Contains(string(settingsBytes), `"colorScheme": "Amber Nocturne"`) {
		t.Fatalf("settings missing colorScheme\nsettings:\n%s", string(settingsBytes))
	}
	if got := stdout.String(); !strings.Contains(strings.ToLower(got), "setup complete") {
		t.Fatalf("stdout missing success message\nstdout:\n%s", got)
	}
}
