# Template Quality Bar

## Purpose

The template system is where the product either feels professional or feels like generated junk.

Every generated file must be clean, minimal, and production-oriented.

## Template rules

1. Generate only what is needed.
2. Avoid giant generic templates.
3. Use comments sparingly.
4. Do not include options the user did not select.
5. Use secure defaults.
6. Avoid exposing internal service ports.
7. Use named volumes for persistent data.
8. Use clear service names.
9. Make files editable by humans.
10. Validate every rendered file.

## Template directory

Suggested structure:

```text
templates/
  compose/
    traefik/
      docker-compose.yml.tmpl
    nginx/
      docker-compose.yml.tmpl
  dockerfile/
    node-next.Dockerfile.tmpl
    node-server.Dockerfile.tmpl
    laravel.Dockerfile.tmpl
    python-fastapi.Dockerfile.tmpl
    python-django.Dockerfile.tmpl
    go.Dockerfile.tmpl
  proxy/
    traefik-dynamic.yml.tmpl
    nginx-default.conf.tmpl
  ci/
    github-actions.yml.tmpl
  scripts/
    deploy.sh.tmpl
    rollback.sh.tmpl
    status.sh.tmpl
    backup.sh.tmpl
  env/
    env.example.tmpl
```

## Template inputs

Templates should use `plan.json`, not raw scanner state.

Scanner → plan → templates.

This keeps generation deterministic.

## Dockerfile quality

A generated Dockerfile should:

- use a stable base image;
- install only needed dependencies;
- use production install where possible;
- use multi-stage build when useful;
- include `.dockerignore`;
- expose the correct internal port;
- run a production command;
- avoid copying secrets;
- prefer non-root user when practical.

## Compose quality

A generated compose file should:

- use explicit service names;
- separate public proxy and internal app network when needed;
- not expose database ports to the host;
- include restart policies;
- include named volumes;
- include environment from `.env`;
- include healthchecks where practical;
- include clear labels for Traefik preset;
- avoid hardcoded secrets.

## Traefik quality

Traefik config should:

- route by domain;
- use HTTPS;
- use internal Docker service port;
- avoid exposing app directly;
- generate only needed labels;
- document required DNS condition.

## Nginx quality

Nginx config should:

- proxy to internal app service;
- set required headers;
- support WebSocket headers;
- include body size only if needed;
- include SSL instructions if cert automation is not fully handled.

## CI quality

CI template should:

- be short;
- use provider secrets;
- SSH into VPS;
- pull latest repo;
- run validation/apply;
- fail on error;
- not contain private keys;
- document required secrets.

## Script quality

Shell scripts should:

- use `set -euo pipefail`;
- print steps;
- fail loudly;
- avoid destructive commands;
- support project root resolution;
- run from expected directory;
- be readable.

## Bad template smell

Avoid templates that:

- include 15 commented alternatives;
- expose every service port;
- use `latest` everywhere;
- use dev commands in production;
- hardcode passwords;
- mount `/` or broad host paths;
- auto-prune Docker;
- assume one framework when confidence is weak.

## Review standard

Before adding a template, ask:

> Would I be comfortable giving this file to a client as my own DevOps work?

If not, do not ship it.
