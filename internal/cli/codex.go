package cli

import (
	"strings"

	"github.com/ychampion/codex-themes/internal/codex"
)

func resolveStatusLineOption(value string) ([]string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return nil, nil
	}
	if strings.Contains(normalized, ",") {
		items := make([]string, 0)
		for _, candidate := range strings.Split(normalized, ",") {
			candidate = strings.TrimSpace(candidate)
			if candidate != "" {
				items = append(items, candidate)
			}
		}
		return items, nil
	}
	return codex.StatusLinePreset(normalized)
}
