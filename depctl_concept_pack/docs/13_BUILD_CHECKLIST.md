# Build Checklist

## Phase 1 — Repository setup

- [ ] Choose final project name.
- [ ] Create GitHub repo.
- [ ] Choose license.
- [ ] Create Go module.
- [ ] Add Cobra CLI.
- [ ] Add config package.
- [ ] Add internal logging/output helper.
- [ ] Add test structure.
- [ ] Add fixtures directory.
- [ ] Add CI for tests.

## Phase 2 — CLI skeleton

Commands:

- [ ] `doctor`
- [ ] `scan`
- [ ] `plan`
- [ ] `write`
- [ ] `validate`
- [ ] `review`
- [ ] `apply`
- [ ] `status`
- [ ] `rollback`

Start with commands printing placeholders.

## Phase 3 — Scanner

- [ ] Implement evidence model.
- [ ] Implement confidence model.
- [ ] Detect repo root.
- [ ] Detect monorepo signals.
- [ ] Detect Node.
- [ ] Detect package manager.
- [ ] Detect Node framework.
- [ ] Detect Laravel.
- [ ] Detect Python.
- [ ] Detect Dockerfile.
- [ ] Parse Dockerfile basics.
- [ ] Detect Compose.
- [ ] Parse Compose services.
- [ ] Extract env keys.
- [ ] Flag sensitive keys.
- [ ] Detect CI provider.
- [ ] Write `detected.json`.
- [ ] Write `scan-report.md`.

## Phase 4 — Planner

- [ ] Define `plan.json` schema.
- [ ] Convert detection into plan.
- [ ] Add preset `compose-traefik`.
- [ ] Add domain input.
- [ ] Select public service.
- [ ] Add env requirement list.
- [ ] Add warnings.
- [ ] Add manual steps.
- [ ] Write `plan-report.md`.

## Phase 5 — Templates

- [ ] Add template renderer.
- [ ] Add template validation for unresolved variables.
- [ ] Node Dockerfile template.
- [ ] Laravel Dockerfile template.
- [ ] Python FastAPI template.
- [ ] Python Django template.
- [ ] Compose + Traefik template.
- [ ] `.dockerignore` template.
- [ ] `.env.example` template.
- [ ] deploy script template.
- [ ] rollback script template.
- [ ] status script template.
- [ ] GitHub Actions template.

## Phase 6 — Writer

- [ ] Create `.deploy/`.
- [ ] Backup existing `.deploy/` files.
- [ ] Render files from plan.
- [ ] Write `.deploy/README.md`.
- [ ] Refuse overwrite without confirmation.
- [ ] Support `--force`.

## Phase 7 — Validator

- [ ] Validate `detected.json`.
- [ ] Validate `plan.json`.
- [ ] Validate Compose syntax.
- [ ] Validate Dockerfile presence.
- [ ] Validate env example.
- [ ] Validate no secret values leaked.
- [ ] Validate database ports not exposed.
- [ ] Validate Traefik labels.
- [ ] Write `validation-report.md`.

## Phase 8 — Apply engine

- [ ] Implement dry-run.
- [ ] Read `plan.json`.
- [ ] Re-run validation.
- [ ] Backup current deployment state.
- [ ] Create Docker networks.
- [ ] Run compose build/up.
- [ ] Check container status.
- [ ] Print logs on failure.
- [ ] Write apply log.
- [ ] Make apply idempotent.

## Phase 9 — Rollback

- [ ] List backups.
- [ ] Restore previous generated kit.
- [ ] Re-run compose.
- [ ] Preserve volumes.
- [ ] Write rollback log.

## Phase 10 — Test and polish

- [ ] Create fixture apps.
- [ ] Add unit tests.
- [ ] Add golden tests.
- [ ] Add integration tests.
- [ ] Test on real VPS.
- [ ] Record demo.
- [ ] Write public README.
- [ ] Publish Docker image.
- [ ] Publish binary releases.
