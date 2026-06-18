# Depctl

Depctl turns a single-app repo into a reviewable deployment setup, then deploys it when the target supports automation.

It keeps the safe flow:

```text
scan -> plan -> write -> validate -> review -> apply
```

For day-to-day use, you usually run one short command.

## Install

From Go:

```bash
go install github.com/AnouarMohamed/Depctl@latest
```

From a cloned repo:

```bash
./install.sh
```

## Fast Start

Prepare a VPS deployment kit without deploying:

```bash
depctl setup --domain app.example.com
```

Deploy to a VPS in one flow:

```bash
depctl deploy --domain app.example.com
```

Deploy to Fly.io:

```bash
depctl deploy --target fly --domain app.example.com
```

Deploy to Vercel:

```bash
depctl deploy --target vercel
```

Useful defaults:

- `--target` defaults to `vps`
- `--preset` defaults to `compose-traefik`
- `--output-dir` defaults to `.deploy`
- `--env-file` defaults to `.env`
- `--region` defaults to `iad` for provider targets

## What It Generates

Depctl writes audit files to `.deploy/`:

- `detected.json`
- `plan.json`
- generated reports
- provider notes and scripts
- backups and rollback state

When needed, it also writes root-level deployment files with backups before overwrite:

- `Dockerfile`
- `.dockerignore`
- `fly.toml`
- `vercel.json`

Secret values from `.env` can be imported into Vercel or Fly during deploy, but they are not written into `plan.json`, reports, or logs.

## Commands

```bash
depctl scan
depctl plan --target vps --domain app.example.com
depctl write
depctl validate
depctl review
depctl apply
depctl status
depctl logs
depctl rollback
```

Shortcut commands:

```bash
depctl setup --domain app.example.com
depctl deploy --target fly --domain app.example.com
```

## Targets

### VPS

The default target. Generates Docker, Docker Compose, reverse proxy config, scripts, and validation reports.

```bash
depctl deploy --domain app.example.com
```

### Fly.io

Requires `flyctl` and either an existing CLI login or `FLY_ACCESS_TOKEN`.

```bash
depctl deploy --target fly --domain app.example.com
```

### Vercel

Best for Next.js and Vite/static apps. Requires Vercel CLI and either an existing CLI login or `VERCEL_TOKEN`.

```bash
depctl deploy --target vercel
```

## Project Status

Depctl currently supports single-app repos first. Monorepos, multi-service deployments, Back4App, Atlas, Render, Railway, DigitalOcean, Kubernetes, and Terraform are intended as later target/provider work.
