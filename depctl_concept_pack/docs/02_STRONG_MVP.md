# Strong MVP

## MVP philosophy

Do not build "supports everything" first.

Build one path that feels production-grade.

The MVP must prove this promise:

> For common web apps on a single VPS, depctl can generate a clean deployment kit that saves real DevOps setup time and avoids common mistakes.

## MVP scope

### Supported host

- Ubuntu/Debian VPS.
- Docker installed or install instructions generated.
- Single server.
- One project at a time.
- Root or sudo user.

### Supported deployment target

Primary preset:

```text
compose-traefik
```

This means:

- Docker Compose;
- Traefik reverse proxy;
- automatic HTTPS via Let's Encrypt;
- one public web service;
- optional internal services;
- deployment through scripts;
- optional GitHub Actions workflow.

Secondary preset, only if easy:

```text
compose-nginx
```

This means:

- Docker Compose;
- Nginx reverse proxy;
- Certbot instructions or generated config;
- one public web service.

Do not include Swarm in MVP unless the Compose path is already excellent.

Swarm becomes v0.2.

### Supported languages/frameworks

MVP should support:

1. Node/Bun apps:
   - package.json;
   - npm/pnpm/yarn/bun;
   - Next.js basic;
   - Vite/static build;
   - Express/Fastify/Nest style server.

2. Laravel/PHP:
   - composer.json;
   - artisan;
   - public directory;
   - php-fpm + nginx container pattern;
   - queue/scheduler detection as warnings first.

3. Python:
   - requirements.txt / pyproject.toml;
   - FastAPI/Uvicorn;
   - Django/Gunicorn;
   - Flask/Gunicorn.

Optional if simple:

4. Go:
   - go.mod;
   - compiled binary;
   - detected port from env or common default.

Do not support Ruby/Java in MVP unless you already use them.

### Supported services

MVP should detect and generate optional service blocks for:

- PostgreSQL;
- MySQL/MariaDB;
- Redis.

MongoDB and object storage can come later.

### Supported CI/CD

MVP should generate:

- GitHub Actions workflow template;
- manual `deploy.sh`;
- `rollback.sh`;
- `status.sh`.

Gitea/GitLab later.

### Supported output

The tool writes only to:

```text
.deploy/
```

MVP output:

```text
.deploy/
  README.md
  plan.json
  detected.json
  docker-compose.yml
  Dockerfile
  .dockerignore
  .env.example
  traefik/
    dynamic.yml
  ci/
    github-actions.yml
  scripts/
    deploy.sh
    rollback.sh
    status.sh
    backup.sh
  reports/
    scan-report.md
    validation-report.md
```

## MVP commands

```bash
depctl doctor
depctl scan
depctl plan --preset compose-traefik --domain app.example.com
depctl write
depctl validate
depctl apply
depctl status
depctl rollback
```

## MVP questions

The tool should ask at most five questions:

1. What domain should serve this app?
2. Which preset? `compose-traefik` or `compose-nginx`.
3. Which service is public if multiple are detected?
4. Which database should be generated if the app seems to need one?
5. Which CI provider? `github`, `none`.

Everything else should be inferred or defaulted.

## MVP anti-scope

Do not build these yet:

- Kubernetes;
- full monitoring;
- full secret manager;
- Portainer automation;
- multi-node Swarm;
- multi-project server dashboard;
- advanced blue/green deployments;
- database migration automation without confirmation;
- automatic DNS management;
- cloud provider APIs;
- random AI-generated config.

## MVP quality bar

The MVP is not done until:

- generated Dockerfiles are usable;
- Compose config runs;
- Traefik routes work;
- HTTPS path is documented;
- `.env.example` is useful;
- scripts are idempotent;
- validation catches common mistakes;
- tool never overwrites without backup;
- reports explain decisions clearly;
- at least 10 real repos have been tested manually.

## MVP demo scenario

A user clones a Node app:

```bash
git clone https://github.com/example/app.git /srv/app
cd /srv/app
depctl plan --preset compose-traefik --domain app.example.com
depctl write
depctl validate
depctl apply
```

The result:

- Dockerfile generated;
- Compose generated;
- Traefik labels configured;
- `.env.example` generated;
- deploy scripts generated;
- GitHub Actions template generated;
- app reachable through the domain.

That is enough to prove the product.
