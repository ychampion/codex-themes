package windows_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/theme"
	ctwindows "github.com/ychampion/codex-themes/internal/windows"
)

func TestApplyThemeUpdatesSchemesAndDefaults(t *testing.T) {
	t.Parallel()

	manifest := mustManifest(t)
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(settingsPath, []byte(`{
  "defaultProfile": "{11111111-1111-1111-1111-111111111111}",
  "profiles": {
    "defaults": {
      "font": { "face": "Cascadia Code" }
    },
    "list": [
      {
        "guid": "{11111111-1111-1111-1111-111111111111}",
        "name": "PowerShell"
      }
    ]
  },
  "schemes": [
    {
      "name": "Existing Theme",
      "background": "#000000",
      "foreground": "#FFFFFF"
    }
  ]
}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := ctwindows.ApplyTheme(settingsPath, manifest, ctwindows.ApplyOptions{})
	if err != nil {
		t.Fatalf("ApplyTheme() error = %v", err)
	}

	bytes, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(bytes)
	for _, fragment := range []string{`"name": "Amber Nocturne"`, `"background": "#12110F"`, `"colorScheme": "Amber Nocturne"`} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("settings missing fragment %q\nsettings:\n%s", fragment, text)
		}
	}
	if result.BackupPath == "" {
		t.Fatal("BackupPath = empty, want path")
	}
	if _, err := os.Stat(result.BackupPath); err != nil {
		t.Fatalf("backup missing: %v", err)
	}
}

func TestDetectSettingsPathFindsStableInstall(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	stable := filepath.Join(base, "Packages", "Microsoft.WindowsTerminal_8wekyb3d8bbwe", "LocalState")
	if err := os.MkdirAll(stable, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	path := filepath.Join(stable, "settings.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := ctwindows.DetectSettingsPath(base)
	if err != nil {
		t.Fatalf("DetectSettingsPath() error = %v", err)
	}
	if got != path {
		t.Fatalf("DetectSettingsPath() = %q, want %q", got, path)
	}
}

func mustManifest(t *testing.T) theme.Manifest {
	t.Helper()
	manifest, err := theme.LoadFile(filepath.Join("..", "..", "themes", "amber-nocturne", "theme.toml"))
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	return manifest
	}

