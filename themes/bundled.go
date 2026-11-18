package themes

import (
	"embed"
	"io/fs"
)

//go:embed */theme.toml
var bundled embed.FS

// Bundled returns the theme manifests shipped with the CLI.
func Bundled() fs.FS {
	return bundled
}
