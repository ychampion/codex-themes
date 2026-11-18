package registry_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ychampion/codex-themes/internal/registry"
)

func TestFindBundledThemeAndLocalOverride(t *testing.T) {
	const themeID = "amber-nocturne"
	bundled, err := registry.Find(themeID, nil)
	if err != nil {
		t.Fatalf("find bundled theme: %v", err)
	}
	if bundled.Name != "Amber Nocturne" {
		t.Fatalf("bundled theme name = %q, want Amber Nocturne", bundled.Name)
	}

	contents, err := os.ReadFile(filepath.Join("..", "..", "themes", themeID, "theme.toml"))
	if err != nil {
		t.Fatal(err)
	}
	override := strings.Replace(string(contents), `name = "Amber Nocturne"`, `name = "Local Amber"`, 1)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "theme.toml"), []byte(override), 0o600); err != nil {
		t.Fatal(err)
	}
	local, err := registry.Find(themeID, []string{dir})
	if err != nil {
		t.Fatalf("find local theme: %v", err)
	}
	if local.Name != "Local Amber" {
		t.Fatalf("local theme name = %q, want Local Amber", local.Name)
	}
}
