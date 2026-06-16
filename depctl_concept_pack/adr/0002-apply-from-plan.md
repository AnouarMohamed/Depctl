# ADR 0002 — Apply From Saved Plan

## Status

Accepted

## Context

If apply rescans the project, the user may review one thing and apply another.

## Decision

`apply` uses `.deploy/plan.json`.

## Consequences

Good:

- predictable;
- auditable;
- safer;
- testable.

Bad:

- user must regenerate plan after changing the project.
