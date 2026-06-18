# Release Process

Depctl releases are tag-driven.

## Requirements

- `main` is green in CI.
- All phase PRs for the release are merged.
- `make verify` passes locally.
- Release notes are written before tagging.

## Create a Release

1. Confirm the next version.

   ```bash
   git fetch origin --tags
   git tag --list 'v*' --sort=-v:refname | head
   ```

2. Update release notes.

   ```bash
   $EDITOR docs/releases/vX.Y.Z.md
   ```

3. Commit the release notes.

4. Create and push the tag.

   ```bash
   git tag -a vX.Y.Z -m "Depctl vX.Y.Z"
   git push origin vX.Y.Z
   ```

5. GitHub Actions will publish:

   - GitHub release notes and binaries
   - Docker image through the existing Docker publish workflow

## Current Release Gate

Run:

```bash
make verify
```

Provider E2E deployments are intentionally manual until credentials and test apps are available.
