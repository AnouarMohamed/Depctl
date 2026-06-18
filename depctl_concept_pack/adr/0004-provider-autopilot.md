# ADR 0004: Provider Autopilot After VPS Core

## Status

Accepted.

## Decision

Depctl remains VPS-first, but the deployment target is now a provider abstraction. The first targets are:

- `vps`
- `vercel`
- `fly`

The safe workflow remains:

```text
scan -> plan -> write -> validate -> review -> apply
```

Shortcut commands can run that workflow for the user:

```bash
depctl setup --domain app.example.com
depctl deploy --target fly --domain app.example.com
depctl deploy --target vercel
```

## Safety Rules

- Provider deploys import `.env` values only during apply/deploy.
- Secret values must not be written to `plan.json`, reports, or logs.
- Root-level provider files are allowed, but existing files must be backed up before overwrite.
- DNS registrar mutation is out of scope for v1.

## Consequences

The planner must produce a v0.2 plan with target, artifacts, checks, credentials, secret imports, actions, and rollback metadata.

VPS remains the quality bar. Vercel and Fly extend the same plan model instead of becoming separate command flows.
