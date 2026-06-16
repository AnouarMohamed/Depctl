# Apply, Rollback, and Idempotency

## Principle

`apply` should do the same thing every time.

If running `apply` twice breaks the server, the design is wrong.

## Apply input

`apply` reads:

```text
.deploy/plan.json
```

It should not rescan the project unless the user explicitly runs:

```bash
depctl scan
depctl plan
depctl write
```

## Apply steps

For Docker Compose MVP:

1. Read `plan.json`.
2. Verify generated files exist.
3. Run validation.
4. Create backup of existing deployment state.
5. Create required Docker networks.
6. Build or pull images.
7. Start/update services with Compose.
8. Wait for containers.
9. Check health.
10. Print routes and status.
11. Write apply log.

## Apply command

```bash
depctl apply
```

Options:

```bash
--yes
--plan .deploy/plan.json
--dry-run
--skip-build
--skip-healthcheck
```

## Dry run

Dry run should print actions without executing them:

```bash
depctl apply --dry-run
```

Example:

```text
Would create network: web
Would run: docker compose -f .deploy/docker-compose.yml up -d --build
Would check route: https://app.example.com
```

## Backups

Before changing anything, store:

```text
.deploy/backups/2026-06-16-1530/
  plan.json
  docker-compose.yml
  apply-log.txt
```

If the tool edits files outside `.deploy/` later, backup them too.

## Rollback

Rollback should:

- list available backups;
- restore previous compose/config files;
- redeploy previous working version;
- never delete database volumes by default.

Command:

```bash
depctl rollback
```

With options:

```bash
--to 2026-06-16-1530
--dry-run
--yes
```

## Idempotent actions

Good actions:

- create network if missing;
- create volume if missing;
- compose up existing services;
- reload proxy safely;
- write same file content if unchanged.

Risky actions:

- delete and recreate database volume;
- overwrite Nginx global config;
- change firewall;
- auto-run migrations;
- prune Docker system.

Risky actions require explicit confirmation or are not MVP.

## Health checks

MVP should check:

- container running;
- port reachable internally;
- Traefik/Nginx route present;
- optional HTTP status;
- Docker logs tail on failure.

## Failure behavior

If apply fails:

- show failed step;
- show command;
- show relevant logs;
- suggest rollback;
- do not hide errors behind generic messages.

## Apply log

Write:

```text
.deploy/reports/apply-log.md
```

Include:

- timestamp;
- plan version;
- commands run;
- results;
- errors;
- rollback hint.

## Migration policy

Do not auto-run database migrations by default.

Generate a script or instruction:

```bash
depctl apply --run-migrations
```

or:

```bash
.deploy/scripts/migrate.sh
```

Reason: migrations can be destructive.

## Data policy

Never delete volumes unless the user explicitly runs a destructive command.

Even rollback should preserve data.
