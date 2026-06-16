# ADR 0003 — MVP Uses Docker Compose + Traefik

## Status

Accepted

## Context

The vision includes Nginx, Docker Swarm, Traefik, Portainer, and CI/CD.

Supporting all of them first will create weak templates and messy behavior.

## Decision

MVP focuses on `compose-traefik`.

## Consequences

Good:

- strong useful path;
- simple VPS deployment;
- HTTPS routing;
- easier testing.

Bad:

- Swarm/Portainer users wait until later.
