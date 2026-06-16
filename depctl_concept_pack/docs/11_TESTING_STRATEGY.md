# Testing Strategy

## Testing goal

depctl must be trusted because deployment tools can break servers.

The scanner, planner, templates, and apply engine all need tests.

## Test pyramid

### Unit tests

Test individual detection rules.

Examples:

- package manager detection;
- framework detection;
- port detection;
- env key extraction;
- Dockerfile parsing;
- compose parsing;
- sensitive key detection.

### Golden file tests

Input fixture repo → expected generated files.

These are essential.

Example:

```text
fixtures/node-next-basic/
  package.json
  expected/
    detected.json
    plan.json
    docker-compose.yml
    Dockerfile
```

### Integration tests

Run generated deployment kit with Docker in CI.

For MVP:

- Node app;
- Laravel app;
- FastAPI app;
- Django app;
- app + PostgreSQL;
- app + Redis.

### Manual VPS tests

Before release, test on real VPS:

- fresh VPS;
- existing Docker installed;
- existing Traefik running;
- existing project with Dockerfile;
- existing compose project.

## Fixture repositories

Create local tiny apps:

```text
fixtures/
  node-express/
  node-next/
  node-vite-static/
  laravel-basic/
  python-fastapi/
  python-django/
  go-basic/
  existing-dockerfile/
  existing-compose/
  monorepo-basic/
```

## Scanner tests

Every detection should produce evidence.

Test that:

- detection value is correct;
- confidence is reasonable;
- evidence references correct files;
- weak detection produces prompt/warning;
- no secret values are copied.

## Template tests

For every template:

- render with sample plan;
- check no unresolved placeholders;
- check expected service names;
- check no secret values;
- check compose syntax if Docker available.

## Apply tests

Apply engine should be tested carefully.

Use dry-run first:

- expected commands;
- expected backup actions;
- expected network creation;
- expected compose command.

Then integration:

- run compose;
- verify app container starts;
- verify healthcheck;
- verify rollback preserves volumes.

## Security tests

Test:

- `.env` with real-looking secret values;
- generated `.env.example` contains only keys, no values;
- database ports not exposed;
- Docker socket not requested in scan/write;
- destructive commands require confirmation.

## Regression tests

Every bug becomes a fixture.

Example:

> Bug: Next.js app with pnpm generated npm command.

Add fixture and test.

## Release checklist

Before every release:

- [ ] Unit tests pass.
- [ ] Golden tests pass.
- [ ] Docker integration tests pass.
- [ ] At least one real VPS smoke test passes.
- [ ] Generated files reviewed manually.
- [ ] No secrets in test output.
- [ ] README updated.
- [ ] Changelog updated.
