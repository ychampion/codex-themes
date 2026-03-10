package preview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ychampion/codex-themes/internal/theme"
)

type Options struct {
	Terminal string
	NoColor  bool
	Plain    bool
	ModelName string
	WorkingDir string
	GitBranch string
	ContextUsed string
}

func Render(manifest theme.Manifest, options Options) (string, error) {
	frame := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(manifest.TokenOr("border", "#555555"))).
		Foreground(lipgloss.Color(manifest.TokenOr("fg_primary", "#FFFFFF"))).
		Background(lipgloss.Color(manifest.TokenOr("bg_primary", "#000000")))

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("fg_muted", "#888888")))
	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color(manifest.TokenOr("status_fg", "#FFFFFF"))).
		Background(lipgloss.Color(manifest.TokenOr("status_bg", "#333333"))).
		Padding(0, 1)
	assistant := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("assistant_msg", "#87CEEB")))
	user := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("user_msg", "#D8B4FE")))
	tool := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("tool_output", "#98FB98")))
	warning := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("warning", "#FBBF24")))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("error", "#F87171")))
	success := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("success", "#4ADE80")))
	code := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("code_inline", manifest.TokenOr("syntax_function", "#60A5FA"))))

	if options.NoColor || options.Plain {
		frame = frame.UnsetBackground()
	}
	modelName := firstNonEmpty(options.ModelName, "o4-mini")
	workingDir := firstNonEmpty(options.WorkingDir, "~/src/codex-themes")
	gitBranch := firstNonEmpty(options.GitBranch, "main")
	contextUsed := firstNonEmpty(options.ContextUsed, "24%")

	body := []string{
		fmt.Sprintf("%s  %s", strings.ToUpper(manifest.Name), muted.Render(manifest.Description)),
		status.Render(fmt.Sprintf("status  model=%s  dir=%s  branch=%s  context=%s", modelName, workingDir, gitBranch, contextUsed)),
		assistant.Render("assistant  I scanned the repo and found 15 bundled themes."),
		user.Render("user       Preview Amber Nocturne and export a Ghostty palette."),
		tool.Render("tool       go test ./..."),
		success.Render("success    Theme validation passed with accessible contrast."),
		warning.Render("warning    ANSI fallback simplifies gradients in basic SSH terminals."),
		errorStyle.Render("error      Config parse failed — restore from backup with codex-theme rollback."),
		"syntax     " + code.Render("const statusAccent = tokens.status_accent"),
	}

	return frame.Render(strings.Join(body, "\n")), nil
	}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
	}

