# ADR 0001 — Product Category

## Status

Accepted

## Context

The product could become a PaaS, an Ansible wrapper, a Dockerfile generator, or a server dashboard.

## Decision

The product will be a repo-aware VPS deployment-kit generator.

## Consequences

Good:

- clear identity;
- easier MVP;
- less infrastructure scope;
- outputs are useful even if user does not apply them automatically.

Bad:

- does not replace a full platform like Coolify;
- requires users to understand generated files at least a little.
