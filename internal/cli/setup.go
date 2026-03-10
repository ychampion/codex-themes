package cli

import (
	"bufio"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ychampion/codex-themes/internal/install"
	"github.com/ychampion/codex-themes/internal/preview"
	"github.com/ychampion/codex-themes/internal/registry"
	"github.com/ychampion/codex-themes/internal/theme"
	ctwindows "github.com/ychampion/codex-themes/internal/windows"
)

func newSetupCommand(themeDirs *[]string, codexHome *string, noColor *bool) *cobra.Command {
	var themeID string
	var settingsPath string
	var skipWindowsTerminal bool
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Windows-first guided setup for Codex Themes",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifests, err := registry.Discover(searchPaths(*themeDirs))
			if err != nil {
				return err
			}
			if len(manifests) == 0 {
				return fmt.Errorf("no themes found in %s", strings.Join(searchPaths(*themeDirs), ", "))
			}
			chosen, err := chooseTheme(cmd, manifests, themeID, *noColor)
			if err != nil {
				return err
			}

			manager := install.NewManager(install.Options{CodexHome: *codexHome})
			state, err := manager.Apply(chosen, install.ApplyOptions{ExportFormats: []install.ExportFormat{install.ExportFormatWindowsTerminal, install.ExportFormatANSI}})
			if err != nil {
				return err
			}

			var windowsMessage string
			if runtime.GOOS == "windows" && !skipWindowsTerminal {
				resolvedPath := settingsPath
				if strings.TrimSpace(resolvedPath) == "" {
					resolvedPath, err = ctwindows.DetectSettingsPath("")
					if err != nil {
						windowsMessage = "Windows Terminal settings not found, so only Codex metadata was updated."
					} else {
						result, err := ctwindows.ApplyTheme(resolvedPath, chosen, ctwindows.ApplyOptions{})
						if err != nil {
							return err
						}
						windowsMessage = fmt.Sprintf("Windows Terminal updated at %s with scheme %s.", result.SettingsPath, result.SchemeName)
					}
				} else {
					result, err := ctwindows.ApplyTheme(resolvedPath, chosen, ctwindows.ApplyOptions{})
					if err != nil {
						return err
					}
					windowsMessage = fmt.Sprintf("Windows Terminal updated at %s with scheme %s.", result.SettingsPath, result.SchemeName)
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "setup complete\n")
			fmt.Fprintf(cmd.OutOrStdout(), "theme: %s (%s)\n", state.ActiveName, state.ActiveThemeID)
			fmt.Fprintf(cmd.OutOrStdout(), "codex home: %s\n", manager.CodexHome())
			if windowsMessage != "" {
				fmt.Fprintln(cmd.OutOrStdout(), windowsMessage)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "next: open a new Windows Terminal tab and run `codex`")
			return nil
		},
	}
	cmd.Flags().StringVar(&themeID, "theme", "", "Theme id to apply without prompting")
	cmd.Flags().StringVar(&settingsPath, "windows-terminal-settings", "", "Explicit Windows Terminal settings.json path")
	cmd.Flags().BoolVar(&skipWindowsTerminal, "skip-windows-terminal", false, "Skip Windows Terminal integration and only apply Codex metadata")
	return cmd
}

func chooseTheme(cmd *cobra.Command, manifests []theme.Manifest, requested string, noColor bool) (theme.Manifest, error) {
	ordered := prioritizeThemes(manifests)
	if strings.TrimSpace(requested) != "" {
		for _, manifest := range ordered {
			if manifest.ID == requested || strings.EqualFold(manifest.Name, requested) {
				return manifest, nil
			}
		}
		return theme.Manifest{}, fmt.Errorf("theme not found: %s", requested)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Codex Themes setup")
	fmt.Fprintln(cmd.OutOrStdout(), "Choose a theme for Codex metadata + Windows Terminal.")
	for index, manifest := range ordered {
		fmt.Fprintf(cmd.OutOrStdout(), "%2d) %s [%s, %s]\n", index+1, manifest.Name, manifest.Appearance, manifest.Family)
	}
	reader := bufio.NewReader(cmd.InOrStdin())
	defaultChoice := 1
	for {
		fmt.Fprintf(cmd.OutOrStdout(), "Theme number [%d]: ", defaultChoice)
		line, err := reader.ReadString('\n')
		if err != nil && strings.TrimSpace(line) == "" {
			return ordered[defaultChoice-1], nil
		}
		line = strings.TrimSpace(line)
		if line == "" {
			ok, err := confirmTheme(cmd, ordered[defaultChoice-1], noColor, reader)
			if err != nil {
				return theme.Manifest{}, err
			}
			if ok {
				return ordered[defaultChoice-1], nil
			}
			continue
		}
		choice, err := strconv.Atoi(line)
		if err != nil || choice < 1 || choice > len(ordered) {
			fmt.Fprintln(cmd.OutOrStdout(), "Please enter a valid number.")
			continue
		}
		ok, err := confirmTheme(cmd, ordered[choice-1], noColor, reader)
		if err != nil {
			return theme.Manifest{}, err
		}
		if ok {
			return ordered[choice-1], nil
		}
	}
}

func confirmTheme(cmd *cobra.Command, manifest theme.Manifest, noColor bool, reader *bufio.Reader) (bool, error) {
	previewText, err := preview.Render(manifest, preview.Options{NoColor: noColor})
	if err == nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintln(cmd.OutOrStdout(), previewText)
		fmt.Fprintln(cmd.OutOrStdout())
	}
	for {
		fmt.Fprintf(cmd.OutOrStdout(), "Apply %s? [Y/n]: ", manifest.Name)
		line, err := reader.ReadString('\n')
		if err != nil && strings.TrimSpace(line) == "" {
			return true, nil
		}
		normalized := strings.ToLower(strings.TrimSpace(line))
		switch normalized {
		case "", "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Fprintln(cmd.OutOrStdout(), "Please answer y or n.")
		}
	}
}

func prioritizeThemes(manifests []theme.Manifest) []theme.Manifest {
	ordered := make([]theme.Manifest, 0, len(manifests))
	var preferred []theme.Manifest
	for _, manifest := range manifests {
		if manifest.ID == "amber-nocturne" {
			preferred = append(preferred, manifest)
			continue
		}
		ordered = append(ordered, manifest)
	}
	return append(preferred, ordered...)
}
