# Agent Phase Review Checklist

## Purpose

When an agent finishes a phase, a reviewer agent uses this document to
decide whether the phase is complete and safe to hand off.

Each phase has its own checklist. The reviewer must go through every item.
If any item fails, the phase is not complete. Log what failed and return
to the builder agent with specific feedback.

A phase is only complete when every item on its checklist passes.

---

## Phase 1 — Repository setup

The reviewer checks:

**Structure**
- [x] `go.mod` exists with the correct module name
- [x] Go version is 1.22 or later
- [x] All directories from `17_GO_PROJECT_STRUCTURE.md` exist
- [x] `main.go` exists and calls `cmd.Execute()` only
- [x] No business logic in `main.go`

**Dependencies**
- [x] Only approved dependencies from `18_GO_DEPENDENCIES.md` are in `go.mod`
- [x] `go.sum` is present and consistent
- [x] No unlisted dependencies added

**CI**
- [x] GitHub Actions workflow exists for running tests
- [x] `go build ./...` passes
- [x] `go vet ./...` passes

**Types**
- [x] `internal/types/detection.go` matches `19_GO_CORE_TYPES.md` exactly
- [x] `internal/types/plan.go` matches `19_GO_CORE_TYPES.md` exactly
- [x] No fields added or removed

**Output helper**
- [x] `internal/output/output.go` exists
- [x] Has functions for: Info, Success, Warning, Error, Step, Quiet mode
- [x] No other package calls `fmt.Println` directly

---

## Phase 2 — CLI skeleton

The reviewer checks:

**Commands**
- [x] All commands from `15_COMMAND_SET.md` are registered with Cobra
- [x] Each command file is in `cmd/`
- [x] Each command prints a placeholder message and exits cleanly
- [x] No command contains business logic

**Flags**
- [x] `scan` has: `--output-dir`, `--quiet`
- [x] `plan` has: `--preset`, `--domain`, `--ci`, `--output-dir`, `--quiet`
- [x] `write` has: `--output-dir`, `--force`, `--quiet`
- [x] `validate` has: `--output-dir`, `--quiet`
- [x] `apply` has: `--yes`, `--dry-run`, `--skip-build`, `--skip-healthcheck`, `--quiet`
- [x] `rollback` has: `--to`, `--dry-run`, `--yes`
- [x] All flags have help text

**Help output**
- [x] `depctl --help` shows all commands
- [x] `depctl scan --help` shows flags
- [x] `depctl apply --help` shows flags including `--dry-run`

**Build**
- [x] `go build -o depctl .` produces a working binary
- [x] Binary runs on Linux amd64

---

## Phase 3 — Scanner

The reviewer checks:

**Output files**
- [x] `depctl scan` produces `.deploy/detected.json`
- [x] `depctl scan` produces `.deploy/reports/scan-report.md`
- [x] `detected.json` is valid JSON
- [x] `detected.json` unmarshals into `types.Detection` without error
- [x] All required fields are present (no missing keys)

**Detection correctness — run against each fixture**
- [x] `fixtures/node-express/` → runtime.name = "node", framework = "express"
- [x] `fixtures/node-next/` → runtime.name = "node", framework = "nextjs"
- [x] `fixtures/laravel-basic/` → runtime.name = "laravel"
- [x] `fixtures/python-fastapi/` → runtime.name = "python", framework = "fastapi"
- [x] `fixtures/python-django/` → runtime.name = "python", framework = "django"
- [x] `fixtures/go-basic/` → runtime.name = "go"
- [x] `fixtures/existing-dockerfile/` → container.dockerfilePresent = true
- [x] `fixtures/existing-compose/` → container.composePresent = true

**Confidence scoring**
- [x] Scoring follows `21_SCORING_RUBRIC.md` exactly
- [x] Confidence field is one of: "strong", "likely", "weak", "unknown"
- [x] No raw float comparisons in business logic

**Security**
- [x] No `.env` values copied into `detected.json`
- [x] Only env var names appear in `env.keys`
- [x] Sensitive keys are flagged per `08_SECURITY_AND_SECRETS.md`

**Tests**
- [x] Unit tests exist for every detection function
- [x] Golden file tests exist for every fixture
- [x] All tests pass with `go test ./internal/scanner/...`
- [x] No test relies on network access

**Code quality**
- [x] No package other than `internal/scanner/` contains detection logic
- [x] Scanner never writes outside `.deploy/`
- [x] Scanner never modifies project files
- [x] `go vet ./internal/scanner/...` passes

---

## Phase 4 — Planner

The reviewer checks:

**Output files**
- [x] `depctl plan --preset compose-traefik --domain example.com` produces `.deploy/plan.json`
- [x] `depctl plan` produces `.deploy/reports/plan-report.md`
- [x] `plan.json` is valid JSON
- [x] `plan.json` unmarshals into `types.Plan` without error

**Plan correctness**
- [x] `preset` matches the flag value
- [x] `domain` matches the flag value
- [x] `publicService` is set to the detected app service name
- [x] `services` includes app service with correct internalPort
- [x] `services` includes database service if postgres/mysql was detected as likely
- [x] `env.required` includes all keys from `detected.json env.keys`
- [x] `env.sensitive` matches `detected.json env.sensitive`
- [x] `generatedFiles` lists every file the writer will produce
- [x] `actions` lists every action apply will take
- [x] `warnings` carries forward scanner warnings plus any new planner warnings
- [x] `manualSteps` includes DNS and secret instructions

**Reads from detection**
- [x] Planner reads `.deploy/detected.json`, not the project directory
- [x] Planner does not re-scan any files
- [x] If `detected.json` is missing, planner errors with a clear message

**apply safety**
- [x] `plan.json` contains `fileHashes` field (even if empty in phase 4)
- [x] Planner does not run any Docker commands

**Tests**
- [x] Unit tests exist for plan generation from sample detections
- [x] Tests cover: Node+Traefik, Laravel+Traefik, Python+Traefik
- [x] Tests cover: with database, without database
- [x] All tests pass

---

## Phase 5 — Templates and Writer

The reviewer checks:

**Templates present**
- [x] All templates listed in `10_TEMPLATE_QUALITY_BAR.md` exist in `templates/`
- [x] Templates are embedded with `//go:embed`

**Generated file quality — check each against quality bar in `10_TEMPLATE_QUALITY_BAR.md`**

Dockerfile (Node Next.js):
- [x] Uses multi-stage build
- [x] Builder stage: correct base image, correct package manager install, build command
- [x] Runner stage: minimal image, only production files copied
- [x] Runs as non-root user
- [x] EXPOSE matches detected port
- [x] CMD uses production start command
- [x] No dev dependencies in runner stage

docker-compose.yml (Traefik preset):
- [x] Explicit service names
- [x] App service has Traefik labels for domain routing
- [x] App service has Traefik label for HTTPS
- [x] Database port NOT exposed to host
- [x] Named volumes for database
- [x] `restart: unless-stopped` on all services
- [x] `env_file: .env` on app service
- [x] Internal network separates app from proxy network
- [x] No hardcoded secrets

deploy.sh:
- [x] Starts with `#!/bin/bash` and `set -euo pipefail`
- [x] Prints each step
- [x] Runs `docker compose up -d --build`
- [x] Checks container status after deploy
- [x] Idempotent — safe to run twice

.env.example:
- [x] Contains all keys from `plan.env.required`
- [x] All values are empty or safe placeholders
- [x] Sensitive keys have a comment marking them as sensitive
- [x] Does not contain any real values

**Writer behaviour**
- [x] Writer reads only from `plan.json`
- [x] Writer never re-scans the project
- [x] Writer backs up existing `.deploy/` files before overwriting
- [x] Writer refuses to overwrite without `--force` if files exist
- [x] Writer generates `.deploy/.gitignore`
- [x] Writer generates `.deploy/README.md`

**No unresolved placeholders**
- [x] All generated files have no `{{` or `}}` remaining after rendering
- [x] Validator catches unresolved placeholders

**Tests**
- [x] Each template is tested with a sample plan
- [x] Golden file tests compare output against `testdata/`
- [x] Tests check no unresolved placeholders
- [x] Tests check no secret values in output

---

## Phase 6 — Validator

The reviewer checks:

**Validation runs cleanly**
- [x] `depctl validate` on a good kit exits 0
- [x] `depctl validate` on a bad kit exits non-zero
- [x] `depctl validate` produces `.deploy/reports/validation-report.md`

**Checks implemented**
- [x] Compose file parses as valid YAML
- [x] Dockerfile exists if plan requires it
- [x] No database ports exposed to host
- [x] No unresolved template placeholders in any generated file
- [x] `.env.example` contains all keys from `plan.env.required`
- [x] Traefik labels are present and correctly formatted
- [x] `plan.json` exists and parses cleanly
- [x] Warns if runtime confidence is "weak"
- [x] Errors if domain is empty

**Report quality**
- [x] Report shows passed checks, warnings, and blocking errors separately
- [x] Report states clearly whether `depctl apply` is allowed
- [x] Report suggests fixes for blocking errors

---

## Phase 7 — Apply Engine

The reviewer checks:

**Dry run**
- [x] `depctl apply --dry-run` prints all actions without executing any
- [x] Dry run output shows exact Docker commands that would run
- [x] Dry run exits 0

**Apply behaviour**
- [x] Reads `.deploy/plan.json` — does not rescan
- [x] Re-runs validation before applying — blocks if validation fails
- [x] Creates Docker network if missing
- [x] Runs `docker compose up -d --build`
- [x] Waits for containers to be healthy
- [x] Writes backup to `.deploy/backups/<timestamp>/`

**Idempotency**
- [x] Running apply twice does not break the server
- [x] Running apply twice does not delete and recreate volumes
- [x] Running apply on an already-running deployment updates it cleanly

**Failure handling**
- [x] On failure: shows which step failed
- [x] On failure: shows relevant Docker logs
- [x] On failure: suggests `depctl rollback`
- [x] Does not hide errors behind generic messages

**Rollback**
- [x] `depctl rollback` lists available backups
- [x] `depctl rollback --to <timestamp>` restores that backup
- [x] Rollback restores compose files and re-runs compose
- [x] Rollback never deletes database volumes
- [x] `depctl rollback --dry-run` shows what would be restored

**Security**
- [x] Apply requires confirmation unless `--yes` is passed
- [x] Destructive operations require explicit flags

---

## Phase 8 — Doctor

The reviewer checks:

- [x] `depctl doctor` checks Docker is installed
- [x] `depctl doctor` checks Docker Compose is installed
- [x] `depctl doctor` checks ports 80 and 443 are not already bound
- [x] `depctl doctor` checks current user can run Docker commands
- [x] `depctl doctor` warns if an existing Traefik container is running
- [x] Output is clear: pass/warn/fail per check
- [x] Doctor never modifies anything

---

## Phase 9 — Integration and polish

The reviewer checks:

**End-to-end test**
- [x] Full flow works: scan → plan → write → validate → apply --dry-run
- [x] Full flow tested against `fixtures/node-next/` on a real or CI machine
- [x] Full flow tested against `fixtures/python-fastapi/`

**Binary**
- [x] `go build -o depctl .` produces a working binary
- [x] Binary runs on Linux amd64 without any dependencies installed
- [x] Binary size is reasonable (under 50MB)

**README**
- [x] Public README explains what depctl is
- [x] README shows the basic flow
- [x] README shows install instructions
- [x] README shows Docker usage

**Release**
- [x] Docker image builds and runs
- [x] Binary release works

---

## Reviewer rules

1. Go through every item. Do not skip items because they seem obvious.
2. Run the commands yourself. Do not trust the builder's description.
3. Check the actual output files, not the code that generates them.
4. If an item is not applicable for this phase, mark it N/A with a reason.
5. Return specific line-level feedback, not vague complaints.
6. A phase passes only when every applicable item is checked.
