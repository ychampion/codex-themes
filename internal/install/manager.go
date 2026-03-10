package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/ychampion/codex-themes/internal/export"
	"github.com/ychampion/codex-themes/internal/theme"
)

type ExportFormat = export.Format

const (
	ExportFormatWindowsTerminal = export.FormatWindowsTerminal
	ExportFormatITerm2          = export.FormatITerm2
	ExportFormatGhostty         = export.FormatGhostty
	ExportFormatAlacritty       = export.FormatAlacritty
	ExportFormatANSI            = export.FormatANSI
)

type Options struct {
	CodexHome string
}

type ApplyOptions struct {
	ExportFormats []ExportFormat
	BackupDir     string
	DryRun        bool
}

type Manager struct {
	codexHome string
}

type State struct {
	ActiveThemeID string         `json:"active_theme_id"`
	ActiveName    string         `json:"active_name"`
	Appearance    string         `json:"appearance"`
	AppliedAt     time.Time      `json:"applied_at"`
	Backups       []string       `json:"backups"`
	Exports       []ExportRecord `json:"exports"`
}

type ExportRecord struct {
	Format ExportFormat `json:"format"`
	Path   string       `json:"path"`
}

func NewManager(options Options) *Manager {
	home := options.CodexHome
	if home == "" {
		home = defaultCodexHome()
	}
	return &Manager{codexHome: home}
}

func (m *Manager) CodexHome() string {
	return m.codexHome
}

func (m *Manager) Apply(manifest theme.Manifest, options ApplyOptions) (State, error) {
	if err := os.MkdirAll(m.codexHome, 0o755); err != nil {
		return State{}, err
	}
	state := State{
		ActiveThemeID: manifest.ID,
		ActiveName:    manifest.Name,
		Appearance:    manifest.Appearance,
		AppliedAt:     time.Now().UTC(),
	}
	if backup, err := m.backupConfig(options.BackupDir); err != nil {
		return State{}, err
	} else if backup != "" {
		state.Backups = append(state.Backups, backup)
	}
	config, err := m.readConfigMap()
	if err != nil {
		return State{}, err
	}
	config["codex_themes"] = map[string]any{
		"active":     manifest.ID,
		"name":       manifest.Name,
		"appearance": manifest.Appearance,
		"family":     manifest.Family,
		"version":    manifest.Version,
		"applied_at": state.AppliedAt.Format(time.RFC3339),
	}
	if !options.DryRun {
		if err := m.writeConfigMap(config); err != nil {
			return State{}, err
		}
	}
	if len(options.ExportFormats) > 0 {
		records, err := m.writeExports(manifest, options.ExportFormats, options.DryRun)
		if err != nil {
			return State{}, err
		}
		state.Exports = records
	}
	if !options.DryRun {
		if err := m.writeState(state); err != nil {
			return State{}, err
		}
	}
	return state, nil
}

func (m *Manager) Current() (State, error) {
	bytes, err := os.ReadFile(m.statePath())
	if err != nil {
		return State{}, err
	}
	var state State
	if err := json.Unmarshal(bytes, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

func (m *Manager) Rollback() error {
	state, err := m.Current()
	if err != nil {
		return err
	}
	if len(state.Backups) == 0 {
		return errors.New("no backup available for rollback")
	}
	backup := state.Backups[len(state.Backups)-1]
	bytes, err := os.ReadFile(backup)
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath(), bytes, 0o600)
	}

func (m *Manager) statePath() string {
	return filepath.Join(m.codexHome, "codex-themes", "state.json")
}

func (m *Manager) configPath() string {
	return filepath.Join(m.codexHome, "config.toml")
}

func (m *Manager) backupConfig(override string) (string, error) {
	configPath := m.configPath()
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	backupDir := override
	if backupDir == "" {
		backupDir = filepath.Join(m.codexHome, "codex-themes", "backups")
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", err
	}
	backupPath := filepath.Join(backupDir, time.Now().UTC().Format("20060102T150405")+"-config.toml")
	if err := os.WriteFile(backupPath, bytes, 0o600); err != nil {
		return "", err
	}
	return backupPath, nil
	}

func (m *Manager) readConfigMap() (map[string]any, error) {
	configPath := m.configPath()
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, err
	}
	var config map[string]any
	if err := toml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("parse %s: %w", configPath, err)
	}
	if config == nil {
		config = map[string]any{}
	}
	return config, nil
	}

func (m *Manager) writeConfigMap(config map[string]any) error {
	bytes, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath(), bytes, 0o600)
	}

func (m *Manager) writeExports(manifest theme.Manifest, formats []ExportFormat, dryRun bool) ([]ExportRecord, error) {
	exportDir := filepath.Join(m.codexHome, "codex-themes", "exports", manifest.ID)
	if !dryRun {
		if err := os.MkdirAll(exportDir, 0o755); err != nil {
			return nil, err
		}
	}
	records := make([]ExportRecord, 0, len(formats))
	for _, format := range formats {
		content, err := export.Render(manifest, format)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(exportDir, export.Filename(manifest, format))
		if !dryRun {
			if err := os.WriteFile(path, content, 0o600); err != nil {
				return nil, err
			}
		}
		records = append(records, ExportRecord{Format: format, Path: path})
	}
	return records, nil
	}

func (m *Manager) writeState(state State) error {
	stateDir := filepath.Dir(m.statePath())
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.statePath(), bytes, 0o600)
	}

func defaultCodexHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".codex"
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, ".codex")
	}
	return filepath.Join(home, ".codex")
}


