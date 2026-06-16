# Command Set

## Philosophy

depctl should feel like a serious systems tool, not a magic script.

The commands should separate thinking from doing:

```text
scan -> plan -> write -> validate -> review -> apply
```

The most important rule is:

```text
apply reads .deploy/plan.json
apply does not rescan and improvise
```

That keeps the tool predictable. The user reviews the exact deployment plan before it touches the VPS.

## MVP commands

### `depctl doctor`

Checks whether the host is ready.

```bash
depctl doctor
```

Checks Docker, Docker Compose, ports 80/443, current user permissions, disk space, Git, and possible existing reverse proxy containers.

### `depctl scan`

Read-only project analysis.

```bash
depctl scan
```

Writes:

```text
.deploy/detected.json
.deploy/reports/scan-report.md
```

It detects runtime, framework, package manager, start/build commands, exposed port, env vars, databases, Redis, existing Docker files, and CI provider.

### `depctl plan`

Turns the scan result plus user choices into a deployment plan.

```bash
depctl plan --preset compose-traefik --domain app.example.com
```

Writes:

```text
.deploy/plan.json
.deploy/reports/plan-report.md
```

For MVP, presets should be limited to:

```text
compose-traefik
compose-nginx
```

### `depctl write`

Generates the deployment kit from the saved plan.

```bash
depctl write
```

Writes Docker, Compose, reverse proxy, env example, CI, and scripts into `.deploy/`.

### `depctl validate`

Checks whether the generated kit is safe enough to apply.

```bash
depctl validate
```

Writes:

```text
.deploy/reports/validation-report.md
```

Validation should check Compose syntax, Dockerfile basics, env completeness, reverse proxy labels/config, missing volumes, dangerous port exposure, unresolved placeholders, and plan consistency.

### `depctl review`

Shows a human-readable summary before deployment.

```bash
depctl review
```

This is the command that tells the user: what was detected, what will be deployed, what files were generated, what is risky, and what manual secrets are still needed.

### `depctl apply`

Applies the reviewed plan.

```bash
depctl apply
```

It should require confirmation unless `--yes` is passed.

### `depctl status`

Shows deployment status.

```bash
depctl status
```

Shows containers, health, domain route, HTTPS state, and last apply result.

### `depctl logs`

Shows useful logs without forcing the user to remember Docker commands.

```bash
depctl logs
depctl logs app
depctl logs proxy
depctl logs db
```

### `depctl rollback`

Restores the last known-good deployment state.

```bash
depctl rollback
```

It must not delete database volumes by default.

## Shortcut command

### `depctl setup`

Runs the safe preparation flow:

```bash
depctl setup --preset compose-traefik --domain app.example.com
```

Equivalent to:

```text
depctl scan
depctl plan
depctl write
depctl validate
depctl review
```

It must not run `apply` automatically.

## Optional later commands

```bash
depctl explain
depctl diff
depctl clean
depctl preset list
depctl template list
```

For MVP, `review` and generated reports are more important than a separate `explain` command.
