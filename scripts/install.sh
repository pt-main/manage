#!/usr/bin/env bash
set -e

MODE="${1:-now}"
APP_NAME="manage"
SNAP_NAME="manage-cli"
PKG_PATH="github.com/pt-main/manage/cmd/manage"

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

run_now() {
  if command -v snap &> /dev/null; then
    if snap list "$SNAP_NAME" &> /dev/null; then
      echo "🔄 [online] Updating $SNAP_NAME via Snap Store..."
      sudo snap refresh "$SNAP_NAME" --amend
    else
      echo "📦 [online] Installing $SNAP_NAME from Snap Store..."
      sudo snap install "$SNAP_NAME" --classic
    fi
    sudo snap alias "${SNAP_NAME}.${APP_NAME}" "$APP_NAME" || true

  elif command -v go &> /dev/null; then
    echo "🐹 [online] Installing/updating via go install @latest..."
    go install "${PKG_PATH}@latest"
    check_go_path

  else
    echo "❌ Neither 'snap' nor 'go' was found on your system."
    exit 1
  fi
}

run_dev() {
  if [ ! -f "go.mod" ]; then
    echo "❌ Error: 'dev' mode must be run from the repository root (go.mod not found)."
    exit 1
  fi

  if command -v snap &> /dev/null && command -v snapcraft &> /dev/null; then
    echo "🛠️  [dev] Building local snap package..."
    sudo "$(command -v snapcraft)" pack
    SNAP_FILE=$(ls -t ${SNAP_NAME}_*.snap 2>/dev/null | head -n 1)
    if [ -z "$SNAP_FILE" ]; then
      echo "❌ Failed to find built .snap file."
      exit 1
    fi
    echo "📦 [dev] Installing/updating from ./$SNAP_FILE..."
    sudo snap install "./$SNAP_FILE" --dangerous --classic
    sudo snap alias "${SNAP_NAME}.${APP_NAME}" "$APP_NAME" || true

  elif command -v go &> /dev/null; then
    echo "🐹 [dev] Installing/updating from local source via go install..."
    go install ./cmd/manage
    check_go_path

  else
    echo "❌ Neither 'snap + snapcraft' nor 'go' was found for local build."
    exit 1
  fi
}

case "$MODE" in
  now|online|update|"")
    run_now
    ;;
  dev|--dev|local)
    run_dev
    ;;
  *)
    echo "Usage: $0 [now|dev]"
    echo "  now (default) - install or update latest release from Snap Store / Go"
    echo "  dev           - build and install/update from current directory"
    exit 1
    ;;
esac

echo "✅ Done!"
"$APP_NAME" sys version || true
