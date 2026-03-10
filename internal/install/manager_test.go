package install_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/install"
	"github.com/ychampion/codex-themes/internal/theme"
)

func TestApplyWritesManagedConfigStateAndBackup(t *testing.T) {
	t.Parallel()

	manifest := mustInstallManifest(t)
	codexHome := filepath.Join(t.TempDir(), ".codex")
	if err := os.MkdirAll(codexHome, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	configPath := filepath.Join(codexHome, "config.toml")
	if err := os.WriteFile(configPath, []byte("model = \"o3\"\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	manager := install.NewManager(install.Options{CodexHome: codexHome})
	state, err := manager.Apply(manifest, install.ApplyOptions{ExportFormats: []install.ExportFormat{install.ExportFormatGhostty, install.ExportFormatWindowsTerminal}})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	configText := string(configBytes)
	for _, fragment := range []string{"model = ", "[codex_themes]", "active = 'amber-nocturne'", "appearance = 'dark'"} {
		if !strings.Contains(configText, fragment) {
			t.Fatalf("config missing fragment %q\nconfig:\n%s", fragment, configText)
		}
	}

	if state.ActiveThemeID != "amber-nocturne" {
		t.Fatalf("ActiveThemeID = %q, want amber-nocturne", state.ActiveThemeID)
	}
	if len(state.Backups) != 1 {
		t.Fatalf("len(Backups) = %d, want 1", len(state.Backups))
	}
	if len(state.Exports) != 2 {
		t.Fatalf("len(Exports) = %d, want 2", len(state.Exports))
	}

	for _, backup := range state.Backups {
		if _, err := os.Stat(backup); err != nil {
			t.Fatalf("backup %q not found: %v", backup, err)
		}
	}
	for _, exported := range state.Exports {
		if _, err := os.Stat(exported.Path); err != nil {
			t.Fatalf("export %q not found: %v", exported.Path, err)
		}
	}
}

func TestRollbackRestoresPreviousConfig(t *testing.T) {
	t.Parallel()

	manifest := mustInstallManifest(t)
	codexHome := filepath.Join(t.TempDir(), ".codex")
	if err := os.MkdirAll(codexHome, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	configPath := filepath.Join(codexHome, "config.toml")
	original := "model = \"o3\"\n[features]\nundo = true\n"
	if err := os.WriteFile(configPath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	manager := install.NewManager(install.Options{CodexHome: codexHome})
	if _, err := manager.Apply(manifest, install.ApplyOptions{}); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if err := manager.Rollback(); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := string(configBytes); got != original {
		t.Fatalf("config after rollback = %q, want %q", got, original)
	}
}

func mustInstallManifest(t *testing.T) theme.Manifest {
	t.Helper()

	manifest, err := theme.Load([]byte(validInstallManifest()))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	return manifest
	}

func validInstallManifest() string {
	return `schema = "https://codexthemes.dev/schemas/theme.schema.json"
id = "amber-nocturne"
name = "Amber Nocturne"
family = "signature"
version = "1.0.0"
appearance = "dark"

[tokens]
bg_primary = "#12110F"
bg_secondary = "#191714"
bg_tertiary = "#211D18"
bg_overlay = "#090807CC"
fg_primary = "#F6F1E8"
fg_secondary = "#E7DDCF"
fg_muted = "#B8AFA2"
fg_inverse = "#120F0A"
border = "#2F2A23"
border_subtle = "#24201A"
border_focus = "#D6A24A"
accent_primary = "#D6A24A"
accent_secondary = "#C97A2B"
accent_soft = "#2B2318"
assistant_msg = "#7CC7FF"
user_msg = "#E2A7FF"
tool_output = "#9FD36B"
shell_stdout = "#D9D1C6"
shell_stderr = "#F07C6C"
status_bg = "#1A1713"
status_fg = "#F4EDE2"
status_accent = "#D6A24A"
input_bg = "#15120F"
input_fg = "#F7F2EA"
input_placeholder = "#968A7B"
selection_bg = "#2B241B"
selection_fg = "#FAF6EF"
cursor = "#E0B35A"
link = "#71B7FF"
info = "#5DB7F0"
success = "#78C26D"
warning = "#E2A84B"
error = "#EA6D5D"
diff_add = "#78C26D"
diff_remove = "#EA6D5D"
diff_change = "#E2A84B"
code_inline = "#7AB8FF"
syntax_plain = "#F0E8DB"
syntax_comment = "#8C8174"
syntax_keyword = "#FF8FB1"
syntax_string = "#9DD68C"
syntax_number = "#F2C66D"
syntax_function = "#80BFFF"
syntax_type = "#63D2CF"
syntax_variable = "#F0E8DB"
syntax_operator = "#F0AA8A"

[ansi]
black = "#1B1713"
red = "#EA6D5D"
green = "#78C26D"
yellow = "#E2A84B"
blue = "#80BFFF"
magenta = "#D79BE6"
cyan = "#63D2CF"
white = "#E7DDCF"
bright_black = "#6E6357"
bright_red = "#F08B7D"
bright_green = "#95D889"
bright_yellow = "#EBC06D"
bright_blue = "#9ACEFF"
bright_magenta = "#E5B8F0"
bright_cyan = "#86E2DF"
bright_white = "#F8F3EC"
`
}

