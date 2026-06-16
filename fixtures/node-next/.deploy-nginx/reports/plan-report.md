# Deployment Plan Report

- **Preset:** compose-nginx
- **Target Domain:** nginx.example.com
- **Public Service:** web

## Services Setup
### Service: web
- **Type:** app
- **Build Context:** .
- **Internal Port:** 3000
- **Publicly Exposed:** true

## Environment Requirements
- None

## Files to Generate
- .deploy/docker-compose.yml
- .deploy/.env.example
- .deploy/.gitignore
- .deploy/scripts/deploy.sh
- .deploy/scripts/rollback.sh
- .deploy/scripts/status.sh
- .deploy/scripts/backup.sh
- .deploy/README.md
- .deploy/Dockerfile
- .deploy/.dockerignore
- .deploy/nginx/default.conf
- .deploy/ci/github-actions.yml

## Actions to Execute on Apply
1. **compose_up** using .deploy/docker-compose.yml

## Plan Warnings
- No .env.example found.
- No Dockerfile found.

## Required Manual Steps
1. [ ] Create DNS A record for nginx.example.com pointing to this VPS.
2. [ ] Fill real secret values in .env on the VPS.
3. [ ] Review .deploy/docker-compose.yml before applying.

## Next Steps
Run the writer command to compile the deployment files:
```bash
depctl write
```
