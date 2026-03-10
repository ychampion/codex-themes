package codex

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ychampion/codex-themes/internal/theme"
)

const DefaultStatusLinePreset = "balanced"

var statusLinePresets = map[string][]string{
	"branch-only": {
		"git-branch",
	},
	"minimal": {
		"model-name",
		"current-dir",
		"git-branch",
	},
	"balanced": {
		"model-with-reasoning",
		"current-dir",
		"git-branch",
		"context-used",
		"fast-mode",
	},
	"dense": {
		"model-with-reasoning",
		"current-dir",
		"git-branch",
		"context-used",
		"used-tokens",
		"fast-mode",
	},
}

func ThemeName(manifest theme.Manifest) string {
	if value := strings.TrimSpace(manifest.Targets.Codex.ThemeName); value != "" {
		return value
	}
	return "codex-themes-" + manifest.ID
}

func ThemeFilename(manifest theme.Manifest) string {
	return ThemeName(manifest) + ".tmTheme"
}

func ThemePath(codexHome string, manifest theme.Manifest) string {
	return filepath.Join(codexHome, "themes", ThemeFilename(manifest))
}

func KnownStatusLinePresets() []string {
	keys := make([]string, 0, len(statusLinePresets))
	for key := range statusLinePresets {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func StatusLinePreset(name string) ([]string, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		normalized = DefaultStatusLinePreset
	}
	values, ok := statusLinePresets[normalized]
	if !ok {
		return nil, fmt.Errorf("unknown status line preset %q (supported: %s)", name, strings.Join(KnownStatusLinePresets(), ", "))
	}
	return cloneStrings(values), nil
}

func ResolveStatusLine(manifest theme.Manifest, override []string) []string {
	if len(override) > 0 {
		return cloneStrings(override)
	}
	if len(manifest.Targets.Codex.StatusLine) > 0 {
		return cloneStrings(manifest.Targets.Codex.StatusLine)
	}
	values, _ := StatusLinePreset(DefaultStatusLinePreset)
	return values
}

func RenderTMTheme(manifest theme.Manifest) ([]byte, error) {
	name := fmt.Sprintf("Codex Themes / %s", manifest.Name)
	background := manifest.TokenOr("bg_primary", "#000000")
	foreground := manifest.TokenOr("syntax_plain", manifest.TokenOr("fg_primary", "#FFFFFF"))
	selectionBackground := manifest.TokenOr("selection_bg", manifest.TokenOr("accent_soft", "#333333"))
	selectionForeground := manifest.TokenOr("selection_fg", manifest.TokenOr("fg_primary", "#FFFFFF"))
	lineHighlight := manifest.TokenOr("bg_secondary", background)
	findHighlight := manifest.TokenOr("accent_soft", selectionBackground)
	guides := manifest.TokenOr("border_subtle", lineHighlight)
	activeGuide := manifest.TokenOr("border_focus", manifest.TokenOr("accent_primary", foreground))
	invisibles := manifest.TokenOr("fg_muted", foreground)
	cursor := manifest.TokenOr("cursor", manifest.TokenOr("accent_primary", foreground))

	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n")
	buffer.WriteString("<plist version=\"1.0\"><dict>\n")
	writeKeyString(&buffer, "name", name)
	buffer.WriteString("<key>settings</key><array>\n")

	writeSettingsItem(&buffer, "", map[string]string{
		"background":          background,
		"foreground":          foreground,
		"caret":               cursor,
		"selection":           selectionBackground,
		"selectionForeground": selectionForeground,
		"lineHighlight":       lineHighlight,
		"findHighlight":       findHighlight,
		"invisibles":          invisibles,
		"guide":               guides,
		"activeGuide":         activeGuide,
	})

	writeSettingsItem(&buffer, "comment, punctuation.definition.comment", map[string]string{
		"foreground": manifest.TokenOr("syntax_comment", manifest.TokenOr("fg_muted", foreground)),
	})
	writeSettingsItem(&buffer, "keyword, storage, keyword.control", map[string]string{
		"foreground": manifest.TokenOr("syntax_keyword", manifest.TokenOr("accent_secondary", foreground)),
	})
	writeSettingsItem(&buffer, "string, constant.other.symbol, markup.inline.raw.string", map[string]string{
		"foreground": manifest.TokenOr("syntax_string", manifest.TokenOr("success", foreground)),
	})
	writeSettingsItem(&buffer, "constant.numeric, constant.language, constant.character", map[string]string{
		"foreground": manifest.TokenOr("syntax_number", manifest.TokenOr("warning", foreground)),
	})
	writeSettingsItem(&buffer, "entity.name.function, support.function, meta.function-call", map[string]string{
		"foreground": manifest.TokenOr("syntax_function", manifest.TokenOr("assistant_msg", foreground)),
	})
	writeSettingsItem(&buffer, "entity.name.type, support.type, storage.type, support.class", map[string]string{
		"foreground": manifest.TokenOr("syntax_type", manifest.TokenOr("info", foreground)),
	})
	writeSettingsItem(&buffer, "variable, support.variable, meta.definition.variable", map[string]string{
		"foreground": manifest.TokenOr("syntax_variable", foreground),
	})
	writeSettingsItem(&buffer, "keyword.operator, punctuation.separator, punctuation.accessor", map[string]string{
		"foreground": manifest.TokenOr("syntax_operator", manifest.TokenOr("accent_primary", foreground)),
	})
	writeSettingsItem(&buffer, "markup.bold", map[string]string{
		"foreground": manifest.TokenOr("accent_primary", foreground),
		"fontStyle":  "bold",
	})
	writeSettingsItem(&buffer, "markup.italic", map[string]string{
		"foreground": manifest.TokenOr("accent_secondary", foreground),
		"fontStyle":  "italic",
	})
	writeSettingsItem(&buffer, "markup.inserted, diff.inserted", map[string]string{
		"foreground": manifest.TokenOr("diff_add", manifest.TokenOr("success", foreground)),
		"background": blendHex(manifest.TokenOr("diff_add", manifest.TokenOr("success", foreground)), background, 0.16),
	})
	writeSettingsItem(&buffer, "markup.deleted, diff.deleted", map[string]string{
		"foreground": manifest.TokenOr("diff_remove", manifest.TokenOr("error", foreground)),
		"background": blendHex(manifest.TokenOr("diff_remove", manifest.TokenOr("error", foreground)), background, 0.16),
	})
	writeSettingsItem(&buffer, "markup.changed, diff.changed", map[string]string{
		"foreground": manifest.TokenOr("diff_change", manifest.TokenOr("warning", foreground)),
		"background": blendHex(manifest.TokenOr("diff_change", manifest.TokenOr("warning", foreground)), background, 0.16),
	})

	buffer.WriteString("</array>\n</dict></plist>")
	return buffer.Bytes(), nil
}

func writeSettingsItem(buffer *bytes.Buffer, scope string, settings map[string]string) {
	buffer.WriteString("<dict>")
	if strings.TrimSpace(scope) != "" {
		writeKeyString(buffer, "scope", scope)
	}
	buffer.WriteString("<key>settings</key><dict>")
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		writeKeyString(buffer, key, settings[key])
	}
	buffer.WriteString("</dict></dict>\n")
}

func writeKeyString(buffer *bytes.Buffer, key, value string) {
	buffer.WriteString("<key>")
	buffer.WriteString(xmlEscape(key))
	buffer.WriteString("</key><string>")
	buffer.WriteString(xmlEscape(value))
	buffer.WriteString("</string>")
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(value)
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}

func blendHex(top, bottom string, alpha float64) string {
	tr, tg, tb := parseHexColor(top)
	br, bg, bb := parseHexColor(bottom)
	blend := func(topValue, bottomValue uint8) uint8 {
		result := alpha*float64(topValue) + (1-alpha)*float64(bottomValue)
		if result < 0 {
			result = 0
		}
		if result > 255 {
			result = 255
		}
		return uint8(result + 0.5)
	}
	return fmt.Sprintf("#%02X%02X%02X", blend(tr, br), blend(tg, bg), blend(tb, bb))
}

func parseHexColor(value string) (uint8, uint8, uint8) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(trimmed) < 6 {
		return 0, 0, 0
	}
	var red, green, blue uint8
	_, _ = fmt.Sscanf(strings.ToUpper(trimmed[:6]), "%02X%02X%02X", &red, &green, &blue)
	return red, green, blue
}
