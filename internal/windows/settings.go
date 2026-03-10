package windows

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tailscale/hujson"

	"github.com/ychampion/codex-themes/internal/export"
	"github.com/ychampion/codex-themes/internal/theme"
)

type ApplyOptions struct {
	DryRun    bool
	BackupDir string
}

type ApplyResult struct {
	SettingsPath string
	BackupPath   string
	SchemeName   string
}

func DetectSettingsPath(base string) (string, error) {
	if strings.TrimSpace(base) == "" {
		base = os.Getenv("LOCALAPPDATA")
	}
	if strings.TrimSpace(base) == "" {
		return "", errors.New("LOCALAPPDATA is not set")
	}
	candidates := []string{
		filepath.Join(base, "Packages", "Microsoft.WindowsTerminal_8wekyb3d8bbwe", "LocalState", "settings.json"),
		filepath.Join(base, "Packages", "Microsoft.WindowsTerminalPreview_8wekyb3d8bbwe", "LocalState", "settings.json"),
		filepath.Join(base, "Microsoft", "Windows Terminal", "settings.json"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("windows terminal settings not found under %s", base)
}

func ApplyTheme(settingsPath string, manifest theme.Manifest, options ApplyOptions) (ApplyResult, error) {
	if strings.TrimSpace(settingsPath) == "" {
		var err error
		settingsPath, err = DetectSettingsPath("")
		if err != nil {
			return ApplyResult{}, err
		}
	}
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return ApplyResult{}, err
	}

	original, err := os.ReadFile(settingsPath)
	if err != nil && !os.IsNotExist(err) {
		return ApplyResult{}, err
	}
	settings, err := parseSettings(original)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("parse windows terminal settings: %w", err)
	}
	scheme, schemeName, err := schemeMap(manifest)
	if err != nil {
		return ApplyResult{}, err
	}
	mergeScheme(settings, scheme, schemeName)
	applyColorScheme(settings, schemeName)

	result := ApplyResult{SettingsPath: settingsPath, SchemeName: schemeName}
	if !options.DryRun {
		if len(original) > 0 {
			backupPath, err := writeBackup(settingsPath, original, options.BackupDir)
			if err != nil {
				return ApplyResult{}, err
			}
			result.BackupPath = backupPath
		}
		encoded, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			return ApplyResult{}, err
		}
		encoded = append(encoded, '\n')
		if err := os.WriteFile(settingsPath, encoded, 0o600); err != nil {
			return ApplyResult{}, err
		}
	}
	return result, nil
}

func parseSettings(input []byte) (map[string]any, error) {
	if len(strings.TrimSpace(string(input))) == 0 {
		return map[string]any{}, nil
	}
	parsed, err := hujson.Parse(input)
	if err != nil {
		return nil, err
	}
	parsed.Standardize()
	var settings map[string]any
	if err := json.Unmarshal(parsed.Pack(), &settings); err != nil {
		return nil, err
	}
	if settings == nil {
		settings = map[string]any{}
	}
	return settings, nil
}

func schemeMap(manifest theme.Manifest) (map[string]any, string, error) {
	content, err := export.Render(manifest, export.FormatWindowsTerminal)
	if err != nil {
		return nil, "", err
	}
	var scheme map[string]any
	if err := json.Unmarshal(content, &scheme); err != nil {
		return nil, "", err
	}
	name, _ := scheme["name"].(string)
	if strings.TrimSpace(name) == "" {
		name = manifest.Name
	}
	return scheme, name, nil
}

func mergeScheme(settings map[string]any, scheme map[string]any, schemeName string) {
	schemes, ok := settings["schemes"].([]any)
	if !ok {
		schemes = []any{}
	}
	replaced := false
	for index, raw := range schemes {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		name, _ := item["name"].(string)
		if name == schemeName {
			schemes[index] = scheme
			replaced = true
			break
		}
	}
	if !replaced {
		schemes = append(schemes, scheme)
	}
	settings["schemes"] = schemes
}

func applyColorScheme(settings map[string]any, schemeName string) {
	profiles := ensureMap(settings, "profiles")
	defaults := ensureMap(profiles, "defaults")
	defaults["colorScheme"] = schemeName

	defaultProfile, _ := settings["defaultProfile"].(string)
	list, _ := profiles["list"].([]any)
	for _, raw := range list {
		profile, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		guid, _ := profile["guid"].(string)
		if guid == defaultProfile && strings.TrimSpace(guid) != "" {
			profile["colorScheme"] = schemeName
			return
		}
	}
}

func ensureMap(parent map[string]any, key string) map[string]any {
	if existing, ok := parent[key].(map[string]any); ok {
		return existing
	}
	created := map[string]any{}
	parent[key] = created
	return created
}

func writeBackup(settingsPath string, original []byte, overrideDir string) (string, error) {
	backupDir := overrideDir
	if strings.TrimSpace(backupDir) == "" {
		backupDir = filepath.Join(filepath.Dir(settingsPath), "codex-themes-backups")
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", err
	}
	backupPath := filepath.Join(backupDir, time.Now().UTC().Format("20060102T150405")+"-settings.json")
	if err := os.WriteFile(backupPath, original, 0o600); err != nil {
		return "", err
	}
	return backupPath, nil
}

