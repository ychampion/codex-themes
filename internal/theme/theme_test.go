package theme_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/theme"
)

func TestLoadManifestAcceptsValidTheme(t *testing.T) {
	t.Parallel()

	path := writeThemeFile(t, validManifest())

	manifest, err := theme.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}

	if manifest.ID != "amber-nocturne" {
		t.Fatalf("manifest.ID = %q, want amber-nocturne", manifest.ID)
	}

	if manifest.Tokens["status_accent"] != "#D6A24A" {
		t.Fatalf("status_accent = %q, want #D6A24A", manifest.Tokens["status_accent"])
	}
	if manifest.ANSI["bright_blue"] != "#9ACEFF" {
		t.Fatalf("bright_blue = %q, want #9ACEFF", manifest.ANSI["bright_blue"])
	}

	if len(manifest.MissingTokens()) != 0 {
		t.Fatalf("MissingTokens() = %v, want none", manifest.MissingTokens())
	}
	if warnings := manifest.Validate(); len(warnings) != 0 {
		t.Fatalf("Validate() warnings = %v, want none", warnings)
	}
	if !manifest.IsAccessible() {
		t.Fatal("IsAccessible() = false, want true")
	}
	if manifest.DisplayName() != "Amber Nocturne (dark)" {
		t.Fatalf("DisplayName() = %q, want Amber Nocturne (dark)", manifest.DisplayName())
	}
	if manifest.Metadata.ThemeDir != filepath.Dir(path) {
		t.Fatalf("ThemeDir = %q, want %q", manifest.Metadata.ThemeDir, filepath.Dir(path))
	}

	if _, ok := manifest.TokenHex("input_bg"); !ok {
		t.Fatal("TokenHex(input_bg) = missing, want present")
	}
	if _, ok := manifest.AnsiHex("green"); !ok {
		t.Fatal("AnsiHex(green) = missing, want present")
	}
	if manifest.AppearsDark() != true {
		t.Fatal("AppearsDark() = false, want true")
	}
	if manifest.TokenOr("missing", "#010203") != "#010203" {
		t.Fatal("TokenOr() did not return fallback")
	}

	for _, token := range theme.RequiredTokenNames() {
		if _, ok := manifest.Tokens[token]; !ok {
			t.Fatalf("required token %q missing after load", token)
		}
	}
}

func TestLoadManifestRejectsMissingRequiredToken(t *testing.T) {
	t.Parallel()

	path := writeThemeFile(t, strings.Replace(validManifest(), "assistant_msg = \"#7CC7FF\"\n", "", 1))

	_, err := theme.LoadFile(path)
	if err == nil {
		t.Fatal("LoadFile() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "assistant_msg") {
		t.Fatalf("LoadFile() error = %q, want assistant_msg mention", err)
	}
}

func TestLoadManifestRejectsInvalidColor(t *testing.T) {
	t.Parallel()

	path := writeThemeFile(t, strings.Replace(validManifest(), "border_focus = \"#D6A24A\"", "border_focus = \"not-a-color\"", 1))

	_, err := theme.LoadFile(path)
	if err == nil {
		t.Fatal("LoadFile() error = nil, want invalid color error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "border_focus") {
		t.Fatalf("LoadFile() error = %q, want border_focus mention", err)
	}
}

func writeThemeFile(t *testing.T, body string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "theme.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
	}

func validManifest() string {
	return `schema = "https://codexthemes.dev/schemas/theme.schema.json"
id = "amber-nocturne"
name = "Amber Nocturne"
family = "signature"
version = "1.0.0"
appearance = "dark"
description = "Warm premium amber theme designed for long coding sessions."
author = "Codex Themes"
license = "MIT"
homepage = "https://github.com/ychampion/codex-themes"
tags = ["warm", "premium", "dark", "amber"]
codex_min_version = "0.1.0"

[accessibility]
min_contrast_body = 4.5
min_contrast_ui = 3.0
color_blind_safe = true

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

[targets.codex]
config_patch = { mode = "merge", path = "config/codex.toml" }

[targets.windows_terminal]
scheme_name = "Amber Nocturne"

[targets.iterm2]
profile_name = "Amber Nocturne"

[targets.ghostty]
theme_name = "Amber Nocturne"

[targets.alacritty]
theme_name = "Amber Nocturne"
`
}


