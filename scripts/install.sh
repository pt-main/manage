#!/usr/bin/env bash
set -e

MODE="${1:-now}"
APP_NAME="manage"
SNAP_NAME="manage-cli"
PKG_PATH="https://github.com/pt-main/manage/cmd/manage"

check_go_path() {
  local gobin
  gobin=$(go env GOBIN)
  if [ -z "$gobin" ]; then
    gobin="$(go env GOPATH)/bin"
  fi

  if [[ ":$PATH:" != *":$gobin:"* ]]; then
    echo "⚠️  Note: $gobin is not in your PATH."
    echo "   Add this line to your shell config (~/.bashrc, ~/.zshrc, or config.fish):"
    echo "   export PATH=\$PATH:$gobin"
  fi
}

install_now() {
  if command -v snap &> /dev/null; then
    echo "📦 [online] Snap detected. Installing $SNAP_NAME from Snap Store..."
    sudo snap install "$SNAP_NAME"
    sudo snap alias "${SNAP_NAME}.${APP_NAME}" "$APP_NAME" || true
  elif command -v go &> /dev/null; then
    echo "🐹 [online] Go detected. Installing via go install..."
    go install "${PKG_PATH}@latest"
    check_go_path
  else
    echo "❌ Neither 'snap' nor 'go' was found on your system."
    exit 1
  fi
}

install_dev() {
  if [ ! -f "go.mod" ]; then
    echo "❌ Error: 'dev' mode must be run from the repository root (go.mod not found)."
    exit 1
  fi

  if command -v snap &> /dev/null && command -v snapcraft &> /dev/null; then
    echo "🛠️  [dev] Building local snap package..."
    snapcraft pack
    SNAP_FILE=$(ls -t ${SNAP_NAME}_*.snap 2>/dev/null | head -n 1)
    if [ -z "$SNAP_FILE" ]; then
      echo "❌ Failed to find built .snap file."
      exit 1
    fi
    echo "📦 [dev] Installing $SNAP_FILE..."
    sudo snap install "./$SNAP_FILE" --dangerous
    sudo snap alias "${SNAP_NAME}.${APP_NAME}" "$APP_NAME" || true
  elif command -v go &> /dev/null; then
    echo "🐹 [dev] Installing from local source via go install..."
    go install ./cmd/manage
    check_go_path
  else
    echo "❌ Neither 'snap + snapcraft' nor 'go' was found for local build."
    exit 1
  fi
}

case "$MODE" in
  now|online|"")
    install_now
    ;;
  dev|--dev|local)
    install_dev
    ;;
  *)
    echo "Usage: $0 [now|dev]"
    echo "  now (default) - install latest release from Snap Store or Go proxy"
    echo "  dev           - build and install from current directory"
    exit 1
    ;;
esac

echo "✅ Done! Run '$APP_NAME init' to get started."
