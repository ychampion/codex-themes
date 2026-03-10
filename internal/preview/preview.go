package preview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ychampion/codex-themes/internal/codex"
	"github.com/ychampion/codex-themes/internal/theme"
)

type Options struct {
	Terminal    string
	NoColor     bool
	Plain       bool
	ModelName   string
	WorkingDir  string
	GitBranch   string
	ContextUsed string
	StatusLine  []string
}

func Render(manifest theme.Manifest, options Options) (string, error) {
	frame := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(manifest.TokenOr("border", "#555555"))).
		Foreground(lipgloss.Color(manifest.TokenOr("fg_primary", "#FFFFFF"))).
		Background(lipgloss.Color(manifest.TokenOr("bg_primary", "#000000")))

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("fg_muted", "#888888")))
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("status_accent", manifest.TokenOr("accent_primary", "#D6A24A")))).Bold(true)
	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color(manifest.TokenOr("status_fg", "#FFFFFF"))).
		Background(lipgloss.Color(manifest.TokenOr("status_bg", "#333333"))).
		Padding(0, 1)
	assistant := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("assistant_msg", "#87CEEB"))).Bold(true)
	user := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("user_msg", "#D8B4FE"))).Bold(true)
	tool := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("tool_output", "#98FB98"))).Bold(true)
	warning := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("warning", "#FBBF24"))).Bold(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("error", "#F87171"))).Bold(true)
	success := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("success", "#4ADE80"))).Bold(true)
	code := lipgloss.NewStyle().Foreground(lipgloss.Color(manifest.TokenOr("code_inline", manifest.TokenOr("syntax_function", "#60A5FA"))))

	if options.NoColor || options.Plain {
		frame = frame.UnsetBackground()
	}

	modelName := firstNonEmpty(options.ModelName, "o4-mini")
	workingDir := firstNonEmpty(options.WorkingDir, "~/src/codex-themes")
	gitBranch := firstNonEmpty(options.GitBranch, "main")
	contextUsed := firstNonEmpty(options.ContextUsed, "24%")
	statusLine := buildStatusLine(options.StatusLine, modelName, workingDir, gitBranch, contextUsed)

	body := []string{
		fmt.Sprintf("%s  %s", accent.Render("CODEX THEMES"), muted.Render(strings.ToUpper(manifest.Name))),
		status.Render(statusLine),
		renderRail(accent, assistant, muted, "AI", "I scanned the repo and wired a custom Codex theme plus a live status line."),
		renderRail(accent, user, muted, "YOU", "Make the session feel premium, distinctive, and screenshot-worthy."),
		renderRail(accent, tool, muted, "RUN", "codex-theme setup --theme amber-nocturne"),
		renderRail(accent, success, muted, "OK", "Codex syntax theme installed and status line preset applied."),
		renderRail(accent, warning, muted, "WARN", "ANSI fallback trims some accents in low-color SSH terminals."),
		renderRail(accent, errorStyle, muted, "ERR", "Config parse failed — restore instantly with codex-theme rollback."),
		renderRail(accent, accent, muted, "CODE", code.Render("const statusAccent = tokens.status_accent")),
	}

	return frame.Render(strings.Join(body, "\n")), nil
}

func buildStatusLine(statusItems []string, modelName, workingDir, gitBranch, contextUsed string) string {
	if len(statusItems) == 0 {
		statusItems, _ = codex.StatusLinePreset(codex.DefaultStatusLinePreset)
	}
	values := map[string]string{
		"model-name":           "MODEL " + modelName,
		"model-with-reasoning": "MODEL " + modelName + " · REASON high",
		"current-dir":          "DIR " + workingDir,
		"project-root":         "ROOT ~/src",
		"git-branch":           "BRANCH " + gitBranch,
		"context-remaining":    "CTX LEFT 58%",
		"context-used":         "CTX " + contextUsed,
		"five-hour-limit":      "5H 82%",
		"weekly-limit":         "WEEK 74%",
		"codex-version":        "CODEX v0.x",
		"context-window-size":  "WINDOW 200K",
		"used-tokens":          "TOKENS 18K",
		"total-input-tokens":   "IN 12K",
		"total-output-tokens":  "OUT 6K",
		"session-id":           "SESSION 9C4A",
		"fast-mode":            "FAST ON",
	}

	segments := make([]string, 0, len(statusItems))
	for _, item := range statusItems {
		if value := strings.TrimSpace(values[item]); value != "" {
			segments = append(segments, value)
		}
	}
	if len(segments) == 0 {
		segments = []string{values["model-with-reasoning"], values["current-dir"], values["git-branch"], values["context-used"]}
	}
	return strings.Join(segments, "  ◆  ")
}

func renderRail(chrome, labelStyle, bodyStyle lipgloss.Style, label, text string) string {
	return fmt.Sprintf("%s %s %s", chrome.Render("│"), labelStyle.Render(label), bodyStyle.Render(text))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
