# Releasing

## GitHub release

Push a tag like `v0.1.0`.

That triggers `.github/workflows/release.yml`, which:

- runs `go test ./...`
- builds release archives for Windows, macOS, and Linux with GoReleaser
- publishes GitHub release assets
- publishes the npm shim when `NPM_TOKEN` is configured

## Required secrets

- `GITHUB_TOKEN` is provided automatically by GitHub Actions
- `NPM_TOKEN` is required only for `npm publish`

## Windows installer flow

The production Windows one-liner is:

```powershell
irm https://raw.githubusercontent.com/ychampion/codex-themes/main/scripts/install.ps1 | iex
```

For a local repo checkout during development:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Source local -RunSetup
```
