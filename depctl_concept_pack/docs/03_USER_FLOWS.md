# User Flows

## Flow 1 — First deployment on a VPS

### Situation

User has a fresh VPS and cloned repo.

### Steps

```bash
cd /srv/my-app
depctl doctor
depctl scan
depctl plan --preset compose-traefik --domain my-app.com
depctl write
depctl validate
depctl apply
```

### Expected result

- `.deploy/` folder exists.
- Docker image builds.
- Compose stack runs.
- Reverse proxy config exists.
- HTTPS route is ready.
- App is deployed.
- User can inspect generated files.

## Flow 2 — Existing Dockerfile

### Situation

Project already has a Dockerfile.

### Behavior

depctl should:

- read the Dockerfile;
- parse `EXPOSE`, `CMD`, `ENTRYPOINT`, build stages;
- validate it;
- avoid replacing it by default;
- generate warnings if it is risky;
- optionally write a suggested Dockerfile to `.deploy/suggestions/Dockerfile.suggested`.

### Rule

Do not overwrite existing Dockerfiles in MVP.

## Flow 3 — Existing docker-compose.yml

### Situation

Project already has compose config.

### Behavior

depctl should:

- read services, networks, volumes, ports;
- detect likely public service;
- detect databases and Redis;
- avoid replacing it by default;
- generate `.deploy/docker-compose.depctl.yml`;
- explain how it relates to the existing file.

## Flow 4 — Laravel app

### Detection

Signals:

- `composer.json`;
- `artisan`;
- `public/index.php`;
- `.env.example` or `.env`;
- `config/database.php`;
- `QUEUE_CONNECTION`;
- `php artisan`.

### Output

- app container;
- nginx/php-fpm pattern or single container pattern;
- database service if needed;
- Redis if queue/cache suggests it;
- `php artisan key:generate` instructions;
- migration step as manual confirmation;
- storage volume note.

## Flow 5 — Node app

### Detection

Signals:

- `package.json`;
- `scripts.start`;
- `scripts.build`;
- lockfile;
- `next.config.js`, `vite.config.*`, `server.js`, `src/main.*`.

### Output

- multi-stage Dockerfile;
- app service;
- detected port;
- static build or server mode;
- healthcheck where possible;
- package manager selection.

## Flow 6 — Python web app

### Detection

Signals:

- `requirements.txt`;
- `pyproject.toml`;
- `manage.py`;
- `main.py`;
- `app.py`;
- FastAPI/Django/Flask imports;
- `gunicorn`, `uvicorn`.

### Output

- Dockerfile;
- command using gunicorn/uvicorn where detected;
- database service if dependency suggests it;
- collectstatic/migration notes for Django.

## Flow 7 — Review before apply

### Situation

User generated files and wants to inspect them.

### UX

```bash
depctl review
```

The tool should show:

- detected stack;
- confidence;
- generated files;
- warnings;
- required manual inputs;
- apply command.

## Flow 8 — Apply

### Requirements

`apply` must:

- read `.deploy/plan.json`;
- verify hash/signature of generated files if possible;
- confirm target;
- backup existing deployment files;
- create networks;
- run Docker Compose;
- display status;
- produce apply log.

## Flow 9 — Rollback

### Situation

Deployment broke.

### Command

```bash
depctl rollback
```

### Behavior

- restore previous compose/config files;
- redeploy previous version if available;
- show what was restored;
- do not delete data volumes by default.

## Flow 10 — CI/CD generation

### Situation

User wants GitHub Actions.

### Output

`.deploy/ci/github-actions.yml`

### Behavior

Generated workflow should:

- connect to VPS through SSH;
- pull latest code;
- run `depctl apply --plan .deploy/plan.json`;
- avoid storing secrets in repo;
- document required GitHub secrets.
