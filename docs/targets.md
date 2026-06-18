# Depctl Targets

Depctl supports one single-app repo per plan.

## VPS

Default target:

```bash
depctl setup --domain app.example.com
depctl apply
```

One-command deploy:

```bash
depctl deploy --domain app.example.com
```

Writes Docker/Compose files, root `Dockerfile`, `.dockerignore`, scripts, reports, and backups.

## Fly.io

Requires `flyctl` and either `fly auth login` or `FLY_ACCESS_TOKEN`.

```bash
depctl deploy --target fly --domain api.example.com
```

Writes root `Dockerfile`, `.dockerignore`, `fly.toml`, imports `.env` with `fly secrets import --stage`, deploys with an image label, and records rollback state in `.deploy/state/fly.json`.

## Vercel

Requires Vercel CLI and either `vercel login` or `VERCEL_TOKEN`.

```bash
depctl deploy --target vercel
```

Best for Next.js and Vite/static apps. Writes root `vercel.json`, imports `.env` keys into production env, deploys with `vercel --prod --yes`, and records the deployment URL in `.deploy/state/vercel.json`.

## Secrets

Secret values are read from `--env-file` during provider deploy. Values are not stored in `plan.json`, reports, or normal command output.

Default:

```bash
--env-file .env
```

## Rollback

```bash
depctl rollback --to <backup>                 # VPS
depctl rollback --to <deployment-id-or-url>   # Vercel
depctl rollback                               # Fly previous image from state
depctl rollback --to <image-ref>              # Fly explicit image
```
