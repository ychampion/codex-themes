package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ychampion/codex-themes/internal/export"
	"github.com/ychampion/codex-themes/internal/install"
	"github.com/ychampion/codex-themes/internal/preview"
	"github.com/ychampion/codex-themes/internal/registry"
	"github.com/ychampion/codex-themes/internal/theme"
)

type Options struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Version string
	ThemeDirs []string
}

func NewRootCommand(options Options) *cobra.Command {
	stdout := options.Stdout
	stderr := options.Stderr
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	themeDirs := append([]string{}, options.ThemeDirs...)
	var codexHome string
	var noColor bool

	root := &cobra.Command{
		Use:           "codex-theme",
		Short:         "Beautiful terminal identities, one command away.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().StringArrayVar(&themeDirs, "theme-dir", themeDirs, "Additional theme search directory")
	root.PersistentFlags().StringVar(&codexHome, "codex-home", "", "Override Codex config root")
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colorized preview output")
	root.AddCommand(newSetupCommand(&themeDirs, &codexHome, &noColor))

	root.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available themes",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifests, err := registry.Discover(searchPaths(themeDirs))
			if err != nil {
				return err
			}
			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 2, 4, 2, ' ', 0)
			fmt.Fprintln(writer, "ID\tNAME\tAPPEARANCE\tFAMILY")
			for _, manifest := range manifests {
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", manifest.ID, manifest.Name, manifest.Appearance, manifest.Family)
			}
			return writer.Flush()
		},
	})
	root.Commands()[0].Flags().Bool("installed", false, "Show installed/local themes only")

	root.AddCommand(&cobra.Command{
		Use:   "validate <theme.toml>",
		Short: "Validate a theme manifest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := theme.LoadFile(args[0])
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "validated %s (%s)\n", manifest.Name, manifest.ID)
			return nil
		},
	})

	previewCommand := &cobra.Command{
		Use:   "preview <theme>",
		Short: "Render a theme preview",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := registry.Find(args[0], searchPaths(themeDirs))
			if err != nil {
				return err
			}
			output, err := preview.Render(manifest, preview.Options{NoColor: noColor})
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), output)
			return err
		},
	}
	root.AddCommand(previewCommand)

	exportCommand := &cobra.Command{
		Use:   "export-terminal <theme>",
		Short: "Export a terminal palette file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := registry.Find(args[0], searchPaths(themeDirs))
			if err != nil {
				return err
			}
			formatValue, _ := cmd.Flags().GetString("format")
			format, err := export.ParseFormat(formatValue)
			if err != nil {
				return err
			}
			content, err := export.Render(manifest, format)
			if err != nil {
				return err
			}
			outputPath, _ := cmd.Flags().GetString("output")
			writeStdout, _ := cmd.Flags().GetBool("stdout")
			if writeStdout || strings.TrimSpace(outputPath) == "" {
				_, err = cmd.OutOrStdout().Write(content)
				return err
			}
			if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
				return err
			}
			return os.WriteFile(outputPath, content, 0o600)
		},
	}
	exportCommand.Flags().String("format", string(export.FormatGhostty), "Export format")
	exportCommand.Flags().String("output", "", "Output file path")
	exportCommand.Flags().Bool("stdout", false, "Write export to stdout")
	root.AddCommand(exportCommand)

	applyCommand := &cobra.Command{
		Use:   "apply <theme>",
		Short: "Apply a theme and update Codex metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := registry.Find(args[0], searchPaths(themeDirs))
			if err != nil {
				return err
			}
			manager := install.NewManager(install.Options{CodexHome: codexHome})
			exportValues, _ := cmd.Flags().GetStringSlice("export-terminal")
			formats, err := parseFormats(exportValues)
			if err != nil {
				return err
			}
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			state, err := manager.Apply(manifest, install.ApplyOptions{DryRun: dryRun, ExportFormats: formats})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "applied %s (%s) to %s\n", state.ActiveName, state.ActiveThemeID, manager.CodexHome())
			return nil
		},
	}
	applyCommand.Flags().StringSlice("export-terminal", nil, "Export companion terminal files")
	applyCommand.Flags().Bool("dry-run", false, "Print intended changes without writing")
	root.AddCommand(applyCommand)

	root.AddCommand(&cobra.Command{
		Use:   "current",
		Short: "Show the active theme",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := install.NewManager(install.Options{CodexHome: codexHome}).Current()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s) %s\n", state.ActiveName, state.ActiveThemeID, state.Appearance)
			return nil
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "rollback",
		Short: "Restore the last backup config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := install.NewManager(install.Options{CodexHome: codexHome}).Rollback(); err != nil {
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "rollback complete")
			return err
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "doctor",
		Short: "Inspect Codex Themes paths and current state",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := install.NewManager(install.Options{CodexHome: codexHome})
			fmt.Fprintf(cmd.OutOrStdout(), "codex_home\t%s\n", manager.CodexHome())
			if _, err := os.Stat(filepath.Join(manager.CodexHome(), "config.toml")); err == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "config\tok")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "config\tmissing (%v)\n", err)
			}
			if state, err := manager.Current(); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "current\t%s (%s)\n", state.ActiveName, state.ActiveThemeID)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "current\tmissing (%v)\n", err)
			}
			return nil
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "sync",
		Short: "Reapply exports for the active theme",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := install.NewManager(install.Options{CodexHome: codexHome})
			state, err := manager.Current()
			if err != nil {
				return err
			}
			manifest, err := registry.Find(state.ActiveThemeID, searchPaths(themeDirs))
			if err != nil {
				return err
			}
			formats := make([]install.ExportFormat, 0, len(state.Exports))
			for _, record := range state.Exports {
				formats = append(formats, record.Format)
			}
			if len(formats) == 0 {
				formats = append(formats, install.ExportFormatGhostty)
			}
			if _, err := manager.Apply(manifest, install.ApplyOptions{ExportFormats: formats}); err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "sync complete")
			return err
		},
	})

	return root
	}

func searchPaths(custom []string) []string {
	paths := append([]string{}, custom...)
	paths = append(paths, registry.DefaultSearchPaths()...)
	return paths
	}

func parseFormats(values []string) ([]install.ExportFormat, error) {
	var formats []install.ExportFormat
	for _, value := range values {
		for _, candidate := range strings.Split(value, ",") {
			if strings.TrimSpace(candidate) == "" {
				continue
			}
			format, err := export.ParseFormat(candidate)
			if err != nil {
				return nil, err
			}
			formats = append(formats, install.ExportFormat(format))
		}
	}
	return formats, nil
	}

