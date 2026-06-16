# Decision Log

Use this file to avoid changing direction every day.

## Decision 001 — Product category

Decision:

depctl is a repo-aware VPS deployment-kit generator, not a PaaS.

Reason:

The strongest gap is generating readable DevOps files from a repo and applying them safely, not building another dashboard platform.

Status:

Accepted.

## Decision 002 — Review before apply

Decision:

The tool writes files first, then applies only after review.

Reason:

This keeps the tool safe, auditable, and useful for real DevOps workflows.

Status:

Accepted.

## Decision 003 — Apply uses saved plan

Decision:

`apply` reads `.deploy/plan.json` and does not rescan automatically.

Reason:

The user must apply exactly what they reviewed.

Status:

Accepted.

## Decision 004 — MVP preset

Decision:

The main MVP preset is `compose-traefik`.

Reason:

Single VPS + Docker Compose + Traefik covers many real deployments with HTTPS and clean routing.

Status:

Accepted.

## Decision 005 — Swarm timing

Decision:

Docker Swarm is not the first MVP path.

Reason:

Swarm adds complexity. Build one excellent Compose path first.

Status:

Accepted but revisitable.

## Decision 006 — Docker image distribution

Decision:

depctl can be distributed as a Docker image and as a host binary.

Reason:

Docker image is convenient. Host binary is cleaner for apply mode because Docker socket mounting is powerful.

Status:

Accepted.

## Decision 007 — Secret values

Decision:

Never copy secret values from `.env`.

Reason:

Generated files should not leak secrets.

Status:

Accepted.

## Decision 008 — Template quality

Decision:

Generated files must be boring, minimal, and manually editable.

Reason:

The product wins only if outputs feel like real DevOps work, not generated clutter.

Status:

Accepted.
