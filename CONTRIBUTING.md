# Contributing

Thanks for contributing to `Codex Themes`.

## Theme contributions

1. Add a new directory under `themes/<id>/`
2. Create `theme.toml`
3. Add `preview.md`
4. Run:

```bash
go test ./...
codex-theme validate themes/<id>/theme.toml
codex-theme preview <id>
```

## Code contributions

- keep changes focused
- add or update tests before implementation changes when practical
- preserve cross-platform behavior
- prefer semantic tokens over renderer-specific branching


