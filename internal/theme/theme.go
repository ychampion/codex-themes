package theme

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var colorPattern = regexp.MustCompile(`^#(?:[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var requiredTokenNames = []string{
	"bg_primary",
	"bg_secondary",
	"bg_tertiary",
	"bg_overlay",
	"fg_primary",
	"fg_secondary",
	"fg_muted",
	"fg_inverse",
	"border",
	"border_subtle",
	"border_focus",
	"accent_primary",
	"accent_secondary",
	"accent_soft",
	"assistant_msg",
	"user_msg",
	"tool_output",
	"shell_stdout",
	"shell_stderr",
	"status_bg",
	"status_fg",
	"status_accent",
	"input_bg",
	"input_fg",
	"input_placeholder",
	"selection_bg",
	"selection_fg",
	"cursor",
	"link",
	"info",
	"success",
	"warning",
	"error",
	"diff_add",
	"diff_remove",
	"diff_change",
	"code_inline",
	"syntax_plain",
	"syntax_comment",
	"syntax_keyword",
	"syntax_string",
	"syntax_number",
	"syntax_function",
	"syntax_type",
	"syntax_variable",
	"syntax_operator",
}

var requiredANSIKeys = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
	"bright_black",
	"bright_red",
	"bright_green",
	"bright_yellow",
	"bright_blue",
	"bright_magenta",
	"bright_cyan",
	"bright_white",
}

type Manifest struct {
	Schema          string         `toml:"schema" json:"schema"`
	ID              string         `toml:"id" json:"id"`
	Name            string         `toml:"name" json:"name"`
	Family          string         `toml:"family" json:"family"`
	Version         string         `toml:"version" json:"version"`
	Appearance      string         `toml:"appearance" json:"appearance"`
	Description     string         `toml:"description" json:"description"`
	Author          string         `toml:"author" json:"author"`
	License         string         `toml:"license" json:"license"`
	Homepage        string         `toml:"homepage" json:"homepage"`
	Tags            []string       `toml:"tags" json:"tags"`
	CodexMinVersion string         `toml:"codex_min_version" json:"codex_min_version"`
	Accessibility   Accessibility  `toml:"accessibility" json:"accessibility"`
	Tokens          map[string]string `toml:"tokens" json:"tokens"`
	ANSI            map[string]string `toml:"ansi" json:"ansi"`
	Targets         Targets        `toml:"targets" json:"targets"`
	Metadata        Metadata       `toml:"-" json:"metadata"`
}

type Metadata struct {
	Path     string `json:"path"`
	ThemeDir string `json:"theme_dir"`
}

type Accessibility struct {
	MinContrastBody float64 `toml:"min_contrast_body" json:"min_contrast_body"`
	MinContrastUI   float64 `toml:"min_contrast_ui" json:"min_contrast_ui"`
	ColorBlindSafe  bool    `toml:"color_blind_safe" json:"color_blind_safe"`
}

type Targets struct {
	Codex           CodexTarget           `toml:"codex" json:"codex"`
	WindowsTerminal WindowsTerminalTarget `toml:"windows_terminal" json:"windows_terminal"`
	ITerm2          ITerm2Target          `toml:"iterm2" json:"iterm2"`
	Ghostty         GhosttyTarget         `toml:"ghostty" json:"ghostty"`
	Alacritty       AlacrittyTarget       `toml:"alacritty" json:"alacritty"`
}

type CodexTarget struct {
	ConfigPatch ConfigPatch `toml:"config_patch" json:"config_patch"`
}

type ConfigPatch struct {
	Mode string `toml:"mode" json:"mode"`
	Path string `toml:"path" json:"path"`
}

type WindowsTerminalTarget struct {
	SchemeName string `toml:"scheme_name" json:"scheme_name"`
}

type ITerm2Target struct {
	ProfileName string `toml:"profile_name" json:"profile_name"`
}

type GhosttyTarget struct {
	ThemeName string `toml:"theme_name" json:"theme_name"`
}

type AlacrittyTarget struct {
	ThemeName string `toml:"theme_name" json:"theme_name"`
}

func RequiredTokenNames() []string {
	return append([]string(nil), requiredTokenNames...)
}

func RequiredANSIKeys() []string {
	return append([]string(nil), requiredANSIKeys...)
}

func LoadFile(path string) (Manifest, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	manifest, err := Load(bytes)
	if err != nil {
		return Manifest{}, err
	}
	manifest.Metadata.Path = path
	manifest.Metadata.ThemeDir = filepath.Dir(path)
	return manifest, nil
}

func Load(data []byte) (Manifest, error) {
	var manifest Manifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	if manifest.Tokens == nil {
		manifest.Tokens = map[string]string{}
	}
	if manifest.ANSI == nil {
		manifest.ANSI = map[string]string{}
	}
	for name, value := range manifest.Tokens {
		manifest.Tokens[name] = strings.ToUpper(value)
	}
	for name, value := range manifest.ANSI {
		manifest.ANSI[name] = strings.ToUpper(value)
	}
	if err := validateRequiredFields(manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (m Manifest) MissingTokens() []string {
	var missing []string
	for _, token := range requiredTokenNames {
		if strings.TrimSpace(m.Tokens[token]) == "" {
			missing = append(missing, token)
		}
	}
	sort.Strings(missing)
	return missing
}

func (m Manifest) MissingANSI() []string {
	var missing []string
	for _, key := range requiredANSIKeys {
		if strings.TrimSpace(m.ANSI[key]) == "" {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}

func (m Manifest) Validate() []string {
	var warnings []string
	if !m.IsAccessible() {
		warnings = append(warnings, "theme contrast does not meet configured minimums")
	}
	return warnings
}

func (m Manifest) IsAccessible() bool {
	bodyThreshold := m.Accessibility.MinContrastBody
	uiThreshold := m.Accessibility.MinContrastUI
	if bodyThreshold == 0 {
		bodyThreshold = 4.5
	}
	if uiThreshold == 0 {
		uiThreshold = 3.0
	}
	body := contrastRatio(m.TokenOr("fg_primary", "#FFFFFF"), m.TokenOr("bg_primary", "#000000"))
	status := contrastRatio(m.TokenOr("status_fg", m.TokenOr("fg_primary", "#FFFFFF")), m.TokenOr("status_bg", m.TokenOr("bg_secondary", "#000000")))
	input := contrastRatio(m.TokenOr("input_fg", m.TokenOr("fg_primary", "#FFFFFF")), m.TokenOr("input_bg", m.TokenOr("bg_secondary", "#000000")))
	return body >= bodyThreshold && status >= uiThreshold && input >= uiThreshold
}

func (m Manifest) DisplayName() string {
	if m.Appearance == "" {
		return m.Name
	}
	return fmt.Sprintf("%s (%s)", m.Name, m.Appearance)
}

func (m Manifest) AppearsDark() bool {
	if strings.EqualFold(m.Appearance, "dark") {
		return true
	}
	if strings.EqualFold(m.Appearance, "light") {
		return false
	}
	return relativeLuminance(m.TokenOr("bg_primary", "#000000")) < 0.4
}

func (m Manifest) TokenHex(name string) (string, bool) {
	value, ok := m.Tokens[name]
	return value, ok
}

func (m Manifest) AnsiHex(name string) (string, bool) {
	value, ok := m.ANSI[name]
	return value, ok
}

func (m Manifest) TokenOr(name, fallback string) string {
	if value, ok := m.Tokens[name]; ok && value != "" {
		return value
	}
	return fallback
}

func validateRequiredFields(m Manifest) error {
	var problems []string
	for field, value := range map[string]string{
		"id":         m.ID,
		"name":       m.Name,
		"family":     m.Family,
		"version":    m.Version,
		"appearance": m.Appearance,
	} {
		if strings.TrimSpace(value) == "" {
			problems = append(problems, field+" is required")
		}
	}
	if m.Appearance != "dark" && m.Appearance != "light" {
		problems = append(problems, "appearance must be dark or light")
	}
	if missing := m.MissingTokens(); len(missing) > 0 {
		problems = append(problems, "missing required tokens: "+strings.Join(missing, ", "))
	}
	if missing := m.MissingANSI(); len(missing) > 0 {
		problems = append(problems, "missing required ansi colors: "+strings.Join(missing, ", "))
	}
	for name, value := range m.Tokens {
		if err := validateColor(name, value); err != nil {
			problems = append(problems, err.Error())
		}
	}
	for name, value := range m.ANSI {
		if err := validateColor(name, value); err != nil {
			problems = append(problems, err.Error())
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func validateColor(name, value string) error {
	if !colorPattern.MatchString(value) {
		return fmt.Errorf("%s must be a #RRGGBB or #RRGGBBAA color", name)
	}
	trimmed := strings.TrimPrefix(value, "#")
	if _, err := hex.DecodeString(trimmed[:6]); err != nil {
		return fmt.Errorf("%s contains invalid hex", name)
	}
	if len(trimmed) == 8 {
		if _, err := hex.DecodeString(trimmed[6:]); err != nil {
			return fmt.Errorf("%s contains invalid alpha hex", name)
		}
	}
	return nil
}

func contrastRatio(foreground, background string) float64 {
	fg := relativeLuminance(foreground)
	bg := relativeLuminance(background)
	lighter := math.Max(fg, bg)
	darker := math.Min(fg, bg)
	return (lighter + 0.05) / (darker + 0.05)
}

func relativeLuminance(hexColor string) float64 {
	r, g, b := parseRGB(hexColor)
	return 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
}

func parseRGB(hexColor string) (float64, float64, float64) {
	trimmed := strings.TrimPrefix(hexColor, "#")
	if len(trimmed) < 6 {
		return 0, 0, 0
	}
	bytes, err := hex.DecodeString(trimmed[:6])
	if err != nil || len(bytes) != 3 {
		return 0, 0, 0
	}
	return float64(bytes[0]) / 255, float64(bytes[1]) / 255, float64(bytes[2]) / 255
}

func linearize(value float64) float64 {
	if value <= 0.03928 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}


