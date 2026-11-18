package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ychampion/codex-themes/internal/theme"
	"github.com/ychampion/codex-themes/themes"
)

func Discover(paths []string) ([]theme.Manifest, error) {
	seen := map[string]theme.Manifest{}
	for _, root := range normalizePaths(paths) {
		entries, err := findThemeFiles(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, entry := range entries {
			manifest, err := theme.LoadFile(entry)
			if err != nil {
				return nil, fmt.Errorf("load %s: %w", entry, err)
			}
			if _, exists := seen[manifest.ID]; !exists {
				seen[manifest.ID] = manifest
			}
		}
	}
	bundled := themes.Bundled()
	entries, err := fs.Glob(bundled, "*/theme.toml")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		data, err := fs.ReadFile(bundled, entry)
		if err != nil {
			return nil, err
		}
		manifest, err := theme.Load(data)
		if err != nil {
			return nil, fmt.Errorf("load bundled %s: %w", entry, err)
		}
		if _, exists := seen[manifest.ID]; !exists {
			seen[manifest.ID] = manifest
		}
	}
	manifests := make([]theme.Manifest, 0, len(seen))
	for _, manifest := range seen {
		manifests = append(manifests, manifest)
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].Name < manifests[j].Name })
	return manifests, nil
}

func Find(identifier string, paths []string) (theme.Manifest, error) {
	if stat, err := os.Stat(identifier); err == nil {
		if stat.IsDir() {
			return theme.LoadFile(filepath.Join(identifier, "theme.toml"))
		}
		return theme.LoadFile(identifier)
	}
	manifests, err := Discover(paths)
	if err != nil {
		return theme.Manifest{}, err
	}
	needle := strings.ToLower(strings.TrimSpace(identifier))
	for _, manifest := range manifests {
		if strings.ToLower(manifest.ID) == needle || strings.ToLower(manifest.Name) == needle {
			return manifest, nil
		}
	}
	return theme.Manifest{}, errors.New("theme not found: " + identifier)
}

func DefaultSearchPaths() []string {
	var paths []string
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, "themes"))
	}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "themes"))
	}
	if env := os.Getenv("CODEX_THEMES_DIR"); strings.TrimSpace(env) != "" {
		paths = append([]string{env}, paths...)
	}
	return normalizePaths(paths)
}

func normalizePaths(paths []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			abs = path
		}
		if _, exists := seen[abs]; exists {
			continue
		}
		seen[abs] = struct{}{}
		result = append(result, abs)
	}
	return result
}

func findThemeFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root && filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.EqualFold(d.Name(), "theme.toml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
