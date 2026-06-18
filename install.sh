#!/usr/bin/env sh
set -eu

BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
mkdir -p "$BIN_DIR"

echo "Building depctl..."
go build -o "$BIN_DIR/depctl" ./cmd/depctl

if [ -t 1 ] && [ -s HHHQ ]; then
  printf '\n'
  cat HHHQ
  printf '\n\n'
fi

case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    echo "Installed to $BIN_DIR/depctl"
    echo "Add this to your shell profile if depctl is not found:"
    echo "  export PATH=\"$BIN_DIR:\$PATH\""
    exit 0
    ;;
esac

echo "Installed depctl to $BIN_DIR/depctl"
