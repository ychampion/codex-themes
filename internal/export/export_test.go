package export_test

import (
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/export"
	"github.com/ychampion/codex-themes/internal/theme"
)

func TestWindowsTerminalExportIncludesSchemeAndCoreColors(t *testing.T) {
	t.Parallel()

	manifest := mustLoadManifest(t)

	output, err := export.Render(manifest, export.FormatWindowsTerminal)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	text := string(output)
	for _, fragment := range []string{`"name": "Amber Nocturne"`, `"background": "#12110F"`, `"foreground": "#F6F1E8"`, `"cursorColor": "#E0B35A"`, `"selectionBackground": "#2B241B"`} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("Render() missing fragment %q\noutput:\n%s", fragment, text)
		}
	}
}

func TestGhosttyExportUsesExpectedKeys(t *testing.T) {
	t.Parallel()

	manifest := mustLoadManifest(t)

	output, err := export.Render(manifest, export.FormatGhostty)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	text := string(output)
	for _, fragment := range []string{"background = #12110F", "foreground = #F6F1E8", "cursor-color = #E0B35A", "selection-background = #2B241B"} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("Render() missing fragment %q\noutput:\n%s", fragment, text)
		}
	}
	if !strings.Contains(text, "palette = 0=#1B1713") {
		t.Fatalf("Render() missing ANSI palette\noutput:\n%s", text)
	}
}

func mustLoadManifest(t *testing.T) theme.Manifest {
	t.Helper()

	manifest, err := theme.Load([]byte(validManifestForExport()))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	return manifest
	}

func validManifestForExport() string {
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


