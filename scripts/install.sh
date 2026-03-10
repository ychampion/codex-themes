#!/usr/bin/env bash
set -euo pipefail

SOURCE="${SOURCE:-auto}"
VERSION="${VERSION:-latest}"
REPO="${REPO:-ychampion/codex-themes}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
RUN_SETUP="${RUN_SETUP:-0}"

say() {
  printf '[codex-theme] %s\n' "$1"
}

repo_root=""
if [ -n "${BASH_SOURCE[0]:-}" ]; then
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  candidate="$(cd "$script_dir/.." && pwd)"
  if [ -f "$candidate/go.mod" ] && [ -f "$candidate/cmd/codex-theme/main.go" ]; then
    repo_root="$candidate"
  fi
fi

mkdir -p "$INSTALL_DIR"
target_path="$INSTALL_DIR/codex-theme"

install_local() {
  if [ -z "$repo_root" ]; then
    say 'Could not detect repo root for local install.'
    exit 1
  fi

  if command -v go >/dev/null 2>&1; then
    say 'Building local binary with Go.'
    (cd "$repo_root" && go build -o "$target_path" ./cmd/codex-theme)
    chmod +x "$target_path"
    return
  fi

  if [ -f "$repo_root/codex-theme" ]; then
    say 'Copying prebuilt local binary.'
    cp "$repo_root/codex-theme" "$target_path"
    chmod +x "$target_path"
    return
  fi

  say 'Local install requested, but no Go toolchain or prebuilt binary was found.'
  exit 1
}

install_release() {
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) arch="x64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) say "Unsupported architecture: $arch"; exit 1 ;;
  esac

  os="$(uname | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux|darwin) ;;
    *) say "Unsupported OS: $os"; exit 1 ;;
  esac

  if [ "$VERSION" = "latest" ]; then
    say 'Resolving latest GitHub release.'
    tag="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p' | head -n1)"
  else
    tag="$VERSION"
    case "$tag" in
      v*) ;;
      *) tag="v$tag" ;;
    esac
  fi

  archive="codex-theme_${os}_${arch}.tar.gz"
  url="https://github.com/$REPO/releases/download/$tag/$archive"
  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  say "Downloading $url"
  curl -fsSL "$url" -o "$tmp_dir/$archive"
  tar -xzf "$tmp_dir/$archive" -C "$tmp_dir"

  binary_path="$(find "$tmp_dir" -type f -name codex-theme | head -n1)"
  if [ -z "$binary_path" ]; then
    say 'Archive did not contain codex-theme.'
    exit 1
  fi

  cp "$binary_path" "$target_path"
  chmod +x "$target_path"
}

case "$SOURCE" in
  auto)
    if [ -n "$repo_root" ]; then
      mode="local"
    else
      mode="release"
    fi
    ;;
  local|release)
    mode="$SOURCE"
    ;;
  *)
    say "Unsupported SOURCE=$SOURCE"
    exit 1
    ;;
esac

say "Install mode: $mode"
if [ "$mode" = "local" ]; then
  install_local
else
  install_release
fi

say "Installed to $target_path"
say "Make sure $INSTALL_DIR is on your PATH, then run: codex-theme setup"

if [ "$RUN_SETUP" = "1" ]; then
  exec "$target_path" setup
fi
