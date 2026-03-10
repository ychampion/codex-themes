# Codex Themes

Beautiful terminal identities, one command away.

`Codex Themes` is an open-source theme manager for Codex CLI. It validates community themes, previews them in the terminal, exports matching terminal palettes, installs a real custom Codex `.tmTheme`, and safely writes reversible Codex config into `~/.codex/config.toml`.

## What ships today

- Semantic theme manifest format in `theme.toml`
- Go CLI with `setup`, `list`, `preview`, `apply`, `sync`, `doctor`, `current`, `rollback`, `validate`, and `export-terminal`
- Safe config backups and rollback support
- Real Codex TUI integration via `~/.codex/themes/*.tmTheme`
- Branded Codex status line presets with model, directory, branch, context, and fast-mode indicators
- Terminal companion exports for Windows Terminal, iTerm2, Ghostty, Alacritty, and ANSI JSON
- Accessibility-aware theme validation
- Bundled premium theme collection under `themes/`

## Install

### Windows right now from this repo

Your screenshot shows the command is not installed on your PATH yet. From `C:\Users\gamin\codex-themes`, run:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Source local -RunSetup
```

That installs `codex-theme` into your user programs folder, adds it to your user `PATH`, and opens the guided theme picker immediately.

If you want a direct one-command install + apply for a git-branch-only footer:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Source local -Theme tokyo-night -StatusLine branch-only
```

There is also a shorter wrapper from the repo root:

```cmd
.\install.cmd -Source local -Theme tokyo-night -StatusLine branch-only
```

### From source

```bash
go install github.com/ychampion/codex-themes/cmd/codex-theme@latest
```

### Windows one-line installer for published releases

```powershell
irm https://raw.githubusercontent.com/ychampion/codex-themes/main/scripts/install.ps1 | iex
```

### macOS / Linux one-line installer for published releases

```bash
curl -fsSL https://raw.githubusercontent.com/ychampion/codex-themes/main/scripts/install.sh | bash
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
codex-theme apply amber-nocturne --status-line balanced --export-terminal ghostty,alacritty
codex-theme current
codex-theme doctor
```

## One-command setup

The fastest way to try the full experience on any OS is:

```bash
codex-theme setup
```

What it does now:

- opens a guided theme picker
- shows a branded live preview before applying
- writes `[tui]` with a real Codex custom theme + status line config
- writes the managed `[codex_themes]` section into `~/.codex/config.toml`
- installs `~/.codex/themes/codex-themes-<theme>.tmTheme`
- exports companion palette files for common terminals
- patches Windows Terminal automatically when available

### Windows / PowerShell

If you use Codex inside Windows Terminal, this is the smoothest path:

```powershell
codex-theme setup
```

After setup finishes, open a new Windows Terminal tab and run:

```powershell
codex
```

### macOS / Linux

Run the same command in your shell:

```bash
codex-theme setup
```

`Codex Themes` automatically installs the Codex theme and exports companion files for terminals like Ghostty, Alacritty, and iTerm2. Use `codex-theme doctor` to see the exported file paths.

If you want a non-interactive version:

```bash
codex-theme setup --theme amber-nocturne --status-line balanced
```

If you only want the Git branch in the Codex footer:

```bash
codex-theme apply tokyo-night --status-line branch-only
```

If you are still in the repo and do not want to install globally yet, the immediate fallback is:

```powershell
.\codex-theme.exe setup
```

Or from anywhere on Windows, without installing globally first:

```powershell
& "C:\Users\gamin\codex-themes\codex-theme.exe" apply tokyo-night --status-line branch-only
```

## How to verify it in Codex

After `setup` or `apply`, check three things:

1. Run `codex-theme doctor`
   - you should see `codex_theme`, `codex_theme_file`, and `status_line`
2. Start a fresh terminal session and run `codex`
   - the footer should show a richer status line with model, directory, branch, and context usage
3. In Codex, open `/theme`
   - you should see the installed custom theme name like `codex-themes-amber-nocturne`

You can also inspect `~/.codex/config.toml` directly and confirm:

- `[tui]`
- `theme = 'codex-themes-<theme>'`
- `status_line = [...]`

## Theme authoring

- Schema: `schemas/theme.schema.json`
- Guide: `docs/theme-authoring.md`
- Release flow: `docs/releasing.md`
- Validate locally: `codex-theme validate path/to/theme.toml`

## Development

```bash
go test ./...
go build ./cmd/codex-theme
```

## Notes on Codex integration

Current Codex builds already support two important hooks that `Codex Themes` uses today:

- custom `.tmTheme` files under `~/.codex/themes/`
- built-in `[tui].status_line`

That means the current release gives you a real in-Codex upgrade right now: custom syntax colors plus a richer footer/status line, while companion terminal palettes keep the outer terminal aligned.


