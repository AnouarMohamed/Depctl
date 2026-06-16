# Roadmap

## Version 0.0 — Concept lock

Goal: freeze product direction.

Deliverables:

- manifesto;
- product spec;
- MVP scope;
- CLI commands;
- scanner design;
- output contract;
- template quality bar.

Exit criteria:

- clear one-liner;
- clear MVP;
- clear non-goals;
- no confusion about product category.

## Version 0.1 — Local scanner prototype

Goal: scan projects and produce `detected.json`.

Support:

- Node;
- Laravel;
- Python;
- Dockerfile detection;
- Compose detection;
- env key detection;
- CI provider detection.

Commands:

```bash
depctl scan
```

Exit criteria:

- produces useful scan report;
- confidence/evidence model exists;
- no file generation yet;
- scanner tested on fixtures.

## Version 0.2 — Planner and writer

Goal: generate `.deploy/` kit.

Support preset:

```text
compose-traefik
```

Commands:

```bash
depctl plan
depctl write
depctl review
```

Exit criteria:

- `plan.json` generated;
- Dockerfile generated when needed;
- compose generated;
- Traefik labels/config generated;
- `.env.example` generated;
- scripts generated;
- GitHub Actions template generated.

## Version 0.3 — Validation

Goal: validate before apply.

Commands:

```bash
depctl validate
```

Checks:

- schema;
- compose syntax;
- unresolved template variables;
- missing env keys;
- exposed database ports;
- missing Docker;
- missing domain;
- weak confidence.

Exit criteria:

- validation catches common bad output;
- validation report is useful.

## Version 0.4 — Apply

Goal: deploy from reviewed plan.

Commands:

```bash
depctl apply
depctl status
```

Support:

- Docker Compose;
- backup;
- network creation;
- compose up;
- status report;
- logs on failure.

Exit criteria:

- apply is idempotent;
- dry-run works;
- real VPS test passes.

## Version 0.5 — Rollback

Goal: safe recovery.

Commands:

```bash
depctl rollback
```

Support:

- restore previous generated kit;
- redeploy previous compose state;
- preserve volumes.

Exit criteria:

- rollback works on real test;
- volumes not deleted.

## Version 0.6 — More presets

Add:

```text
compose-nginx
swarm-traefik
```

Only after Compose + Traefik is solid.

## Version 0.7 — More CI providers

Add:

- Gitea Actions;
- GitLab CI;
- manual-only mode improvements.

## Version 0.8 — Monorepo support

Support:

- selecting app folder;
- multiple services;
- one public gateway;
- shared databases;
- frontend + backend deployment.

## Version 1.0 — Stable release

Criteria:

- stable CLI;
- tested templates;
- real documentation;
- at least 30 real repo tests;
- clear warnings;
- no dangerous defaults;
- predictable apply/rollback;
- enough polish for public use.

## Future ideas

- Portainer stack output;
- Docker Swarm secrets;
- migrations assistant;
- queue/worker generator;
- cron/scheduler support;
- firewall hints;
- monitoring hooks;
- backup strategy generator;
- server inventory;
- team config file;
- plugin system.
