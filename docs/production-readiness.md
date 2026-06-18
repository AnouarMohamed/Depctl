# Production Readiness

This checklist defines what must be true before a release is called production-ready.

## Required Gates

- CI runs formatting, vet, tests, race tests, build, and fixture smoke tests.
- `depctl setup` works for VPS, Vercel, and Fly fixture paths.
- `depctl apply --dry-run` prints exact provider actions without secret values.
- Generated provider files are backed up before overwrite.
- Secret values from `.env` do not appear in `plan.json`, reports, or command logs.
- Release tags publish binaries and container images.

## Manual E2E Gates

These require real infrastructure credentials.

- Deploy a Node HTTP app to a VPS.
- Deploy a FastAPI app to Fly.io.
- Deploy a Next.js app to Vercel.
- Run status and logs after each deployment.
- Roll back each target.

## Current Boundary

The tool supports one public app per repo. Monorepos, multi-service deployments, and managed database provisioning are future phases.
