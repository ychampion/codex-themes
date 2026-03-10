# Theme Authoring

Each theme lives in `themes/<id>/theme.toml` and may optionally include a short `preview.md`.

## Rules

- Use a lowercase kebab-case `id`
- Set `appearance` to `dark` or `light`
- Fill every semantic token in `[tokens]`
- Fill the full 16-color ANSI palette in `[ansi]`
- Keep body contrast readable; the validator checks this

## Validate locally

```bash
codex-theme validate themes/amber-nocturne/theme.toml
codex-theme preview themes/amber-nocturne/theme.toml
```

## Terminal exports

Terminal companion files are generated from semantic tokens. The manifest may include target names under `[targets.*]`, but the export generator can synthesize files even when those sections are omitted.

## Codex-specific targets

The `[targets.codex]` table can optionally include extra Codex integration hints:

- `theme_name` — override the generated custom Codex theme filename
- `status_line` — ordered status line item ids to write into `[tui].status_line`

Example:

```toml
[targets.codex]
theme_name = "codex-themes-amber-nocturne"
status_line = ["model-with-reasoning", "current-dir", "git-branch", "context-used", "fast-mode"]
config_patch = { mode = "merge", path = "config/codex.toml" }
```


