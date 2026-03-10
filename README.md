# Codex Themes

Beautiful terminal identities, one command away.

`Codex Themes` is an open-source theme manager for Codex CLI. It validates community themes, previews them in the terminal, exports matching terminal palettes, and safely writes reversible Codex theme metadata into `~/.codex/config.toml`.

## What ships today

- Semantic theme manifest format in `theme.toml`
- Go CLI with `list`, `preview`, `apply`, `sync`, `doctor`, `current`, `rollback`, `validate`, and `export-terminal`
- Safe config backups and rollback support
- Terminal companion exports for Windows Terminal, iTerm2, Ghostty, Alacritty, and ANSI JSON
- Accessibility-aware theme validation
- Bundled premium theme collection under `themes/`

## Install

### From source

```bash
go install github.com/ychampion/codex-themes/cmd/codex-theme@latest
```

### Local development

```bash
go build ./cmd/codex-theme
./codex-theme list
```

### Planned release channels

- `npm i -g codex-theme`
- Homebrew formula generated via GoReleaser
- Scoop and Winget manifests
- GitHub release archives for macOS, Linux, and Windows

## Quickstart

```bash
codex-theme list
codex-theme preview amber-nocturne
codex-theme apply amber-nocturne --export-terminal ghostty,alacritty
codex-theme current
codex-theme doctor
```

## Theme authoring

- Schema: `schemas/theme.schema.json`
- Guide: `docs/theme-authoring.md`
- Validate locally: `codex-theme validate path/to/theme.toml`

## Development

```bash
go test ./...
go build ./cmd/codex-theme
```

## Notes on Codex integration

Current Codex builds tolerate unknown keys in `config.toml`, so `Codex Themes` writes a managed `[codex_themes]` section plus exported companion files without trying to guess unsupported upstream theme keys. This keeps installs reversible and forward-compatible while richer upstream hooks mature.


