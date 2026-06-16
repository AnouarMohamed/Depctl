# CLI UX Specification

## Design principles

1. Commands should be explicit.
2. The default mode should be safe.
3. Users should see decisions before files are written.
4. Users should review files before apply.
5. `apply` should use the saved plan, not fresh guesses.

## Command list

### `doctor`

Checks host readiness.

```bash
depctl doctor
```

Checks:

- OS;
- Docker availability;
- Docker Compose availability;
- current user permissions;
- ports 80/443 availability;
- existing Traefik/Nginx containers;
- DNS resolution for provided domain if available.

### `scan`

Read-only project analysis.

```bash
depctl scan
```

Outputs:

```text
.deploy/detected.json
.deploy/reports/scan-report.md
```

No infrastructure files should be generated here unless `--write` is passed.

### `plan`

Build deployment plan.

```bash
depctl plan --preset compose-traefik --domain app.example.com
```

Outputs:

```text
.deploy/plan.json
.deploy/reports/plan-report.md
```

### `write`

Render files from plan.

```bash
depctl write
```

Outputs deployment kit files.

### `validate`

Validate generated kit.

```bash
depctl validate
```

Checks:

- JSON schema;
- Compose syntax;
- Dockerfile basics;
- missing env vars;
- dangerous exposed ports;
- invalid Traefik labels;
- unresolved template placeholders;
- conflicting service names;
- missing volumes for persistent paths.

### `review`

Human-friendly summary.

```bash
depctl review
```

Shows:

- detected framework;
- selected preset;
- generated files;
- warnings;
- manual steps;
- exact apply command.

### `apply`

Apply reviewed plan.

```bash
depctl apply
```

Should require confirmation unless `--yes` is used.

### `status`

Shows deployment status.

```bash
depctl status
```

Should show:

- containers;
- exposed routes;
- health status;
- last apply time;
- logs hint.

### `rollback`

Rollback previous deployment.

```bash
depctl rollback
```

Must not delete database volumes by default.

## Suggested command aliases

Do not add these in MVP unless needed:

```bash
depctl init
depctl generate
depctl up
```

Explicit names are better for safety.

## Example output

```text
depctl scan complete

Detected app:
  Type: Node / Next.js
  Confidence: 91%
  Package manager: pnpm
  Build command: pnpm build
  Start command: pnpm start
  Public port: 3000

Detected dependencies:
  PostgreSQL: likely
  Redis: not detected

Warnings:
  - No .env.example found
  - DATABASE_URL referenced but not documented
  - No health endpoint detected

Next:
  depctl plan --preset compose-traefik --domain your-domain.com
```

## Interactive mode

Interactive mode should exist, but flags must always work.

Good:

```bash
depctl plan
```

Then prompt.

Also good:

```bash
depctl plan --preset compose-traefik --domain app.example.com --ci github
```

Bad:

```bash
depctl magic
```

## Output style

Keep output direct:

- what was detected;
- what was generated;
- what is risky;
- what to do next.

Avoid cute logs and fake excitement.
