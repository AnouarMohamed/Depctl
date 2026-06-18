#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
if [[ -z "$version" ]]; then
  echo "usage: scripts/release.sh vX.Y.Z" >&2
  exit 2
fi

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
  echo "invalid version: $version" >&2
  echo "expected format: vX.Y.Z or vX.Y.Z-prerelease" >&2
  exit 2
fi

branch="$(git branch --show-current)"
if [[ "$branch" != "main" ]]; then
  echo "refusing to release from branch '$branch'; switch to main after merging release PRs" >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "refusing to release with a dirty worktree" >&2
  git status --short
  exit 1
fi

notes="docs/releases/${version}.md"
if [[ ! -f "$notes" ]]; then
  echo "missing release notes: $notes" >&2
  exit 1
fi

git fetch origin --tags
if git rev-parse "$version" >/dev/null 2>&1; then
  echo "tag already exists locally: $version" >&2
  exit 1
fi
if git ls-remote --exit-code --tags origin "refs/tags/$version" >/dev/null 2>&1; then
  echo "tag already exists on origin: $version" >&2
  exit 1
fi

make verify

git tag -a "$version" -F "$notes"
git push origin "$version"

echo "Released $version"
