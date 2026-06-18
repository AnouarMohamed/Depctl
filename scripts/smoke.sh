#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

run_setup_smoke() {
  local fixture="$1"
  shift

  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  cp -R "$ROOT/fixtures/$fixture/." "$tmp/"
  rm -rf "$tmp/.deploy" "$tmp/.deploy-nginx"

  go run "$ROOT" setup "$tmp" --force --output-dir "$tmp/.deploy" "$@"
  test -f "$tmp/.deploy/plan.json"
  go run "$ROOT" apply --output-dir "$tmp/.deploy" --dry-run
}

run_setup_smoke "node-express" --domain app.example.com
run_setup_smoke "node-next" --target vercel
run_setup_smoke "python-fastapi" --target fly --domain api.example.com

echo "Smoke tests passed."
