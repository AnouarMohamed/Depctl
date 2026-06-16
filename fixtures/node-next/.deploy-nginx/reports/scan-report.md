# Scan Report

## Project Context
- **Name:** node-next
- **Path:** /home/anouar/Documents/Projects/Depctl/fixtures/node-next
- **Git Repository:** false

## Runtime & Framework
- **Runtime:** node
- **Framework:** nextjs
- **Confidence:** strong (95%)

### Evidence:
- package.json: package.json exists
- package.json: dependency next detected
- pnpm-lock.yaml: pnpm lockfile detected

## Build & Expose Specs
- **Package Manager:** pnpm
- **Build Command:** pnpm run build
- **Start Command:** pnpm start
- **Internal Service Port:** 3000 (Confidence: likely)

## Service Dependencies
- **mysql:** Not detected
- **postgres:** Not detected
- **redis:** Not detected

## Docker & CI Environments
- **Existing Dockerfile:** false
- **Existing Docker Compose:** false
- **CI Pipelines:** GitHub Actions (false), GitLab CI (false), Gitea (false)

## Warnings
- [ ] No .env.example found.
- [ ] No Dockerfile found.

