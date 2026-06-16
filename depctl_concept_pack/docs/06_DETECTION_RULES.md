# Detection Rules

## Node / Bun

### Signals

Strong:

- `package.json`
- `scripts.start`
- `scripts.build`
- lockfile: `package-lock.json`, `pnpm-lock.yaml`, `yarn.lock`, `bun.lockb`

Framework signals:

- `next` dependency → Next.js
- `vite` dependency or `vite.config.*` → Vite
- `express` dependency → Express
- `fastify` dependency → Fastify
- `@nestjs/core` dependency → NestJS
- `nuxt` dependency → Nuxt

### Port detection

Priority:

1. Dockerfile `EXPOSE`;
2. environment variable `PORT`;
3. framework default;
4. user prompt.

Defaults:

- Next.js: 3000
- Express/Fastify/Nest: 3000
- Vite preview: 4173
- Vite dev: 5173, but dev mode should not be used for production.

### Package manager

Priority:

1. `pnpm-lock.yaml` → pnpm
2. `yarn.lock` → yarn
3. `bun.lockb` or `bun.lock` → bun
4. `package-lock.json` → npm
5. fallback → npm

### Production command

Priority:

1. `scripts.start`;
2. framework known command;
3. ask user.

## Laravel / PHP

### Signals

Strong:

- `composer.json`;
- `artisan`;
- `public/index.php`;
- `bootstrap/app.php`.

### Dependency signals

- `laravel/framework` → Laravel.
- `predis/predis` or Redis env usage → Redis likely.
- `QUEUE_CONNECTION` → queue support.
- `DB_CONNECTION` → database required.

### Required notes

Laravel deployment usually needs:

- `APP_KEY`;
- database connection;
- storage permissions;
- `php artisan migrate` decision;
- `php artisan storage:link`;
- queue worker if queues are used;
- scheduler if scheduled commands exist.

### MVP behavior

Generate warnings for queue and scheduler first.

Do not auto-run migrations without confirmation.

## Python

### Signals

Strong:

- `requirements.txt`;
- `pyproject.toml`;
- `Pipfile`;
- `manage.py`;
- `main.py`;
- `app.py`.

Framework signals:

- `fastapi` dependency → FastAPI;
- `django` dependency or `manage.py` → Django;
- `flask` dependency → Flask.

### Production command

FastAPI:

```bash
uvicorn main:app --host 0.0.0.0 --port 8000
```

But only if `main.py` and `app` are likely.

Django:

```bash
gunicorn projectname.wsgi:application --bind 0.0.0.0:8000
```

Requires detecting project module.

Flask:

```bash
gunicorn app:app --bind 0.0.0.0:8000
```

Only if `app.py` and `app` exist.

### MVP rule

If Python entrypoint is uncertain, ask user.

## Go

### Signals

- `go.mod`;
- `main.go`;
- `cmd/*/main.go`.

### Build

Generate multi-stage Dockerfile:

- build in golang image;
- copy binary to minimal runtime image.

### Port

Usually env-driven. Ask user if not detected.

## Dockerfile quality checks

Warn if:

- no `EXPOSE`;
- no `CMD` or `ENTRYPOINT`;
- runs as root;
- copies entire context without `.dockerignore`;
- uses dev command;
- no multi-stage build for compiled/frontend apps;
- installs unnecessary dev dependencies in production;
- hardcoded secrets.

## Compose quality checks

Warn if:

- public service exposes raw app port unnecessarily;
- database port is exposed to host;
- missing named volumes for database;
- missing restart policy;
- no network separation;
- no healthcheck for database;
- no Traefik/Nginx routing labels/config;
- secrets appear in compose file values.

## Env detection

Scan:

- `.env.example`;
- `.env`;
- common config files;
- `process.env.X`;
- `os.Getenv("X")`;
- `$_ENV["X"]`;
- `env("X")`;
- `settings.py`.

Do not copy values from `.env`.

Only copy keys.

## Sensitive key patterns

Flag keys containing:

```text
SECRET
TOKEN
PASSWORD
PASS
PRIVATE
KEY
API_KEY
DB_PASS
JWT
SESSION
CREDENTIAL
ACCESS_KEY
```

## Database detection

PostgreSQL likely if:

- `DATABASE_URL` starts with postgres pattern in example;
- dependency `pg`, `psycopg`, `asyncpg`, `postgres`;
- Laravel `DB_CONNECTION=pgsql`.

MySQL likely if:

- dependency `mysql`, `mysql2`, `pymysql`;
- Laravel `DB_CONNECTION=mysql`.

Redis likely if:

- dependency `redis`, `ioredis`, `predis`;
- env `REDIS_URL`;
- Laravel `CACHE_STORE=redis`, `QUEUE_CONNECTION=redis`.

## Monorepo detection

Signals:

- root package manager workspaces;
- `apps/*`;
- `packages/*`;
- `services/*`;
- multiple package manifests.

MVP behavior:

- detect monorepo;
- ask user to select app folder;
- do not auto-generate multi-service deployment unless obvious.
