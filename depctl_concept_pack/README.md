# depctl Concept Pack

**Working name:** depctl  
**Product category:** repo-aware VPS DevOps scaffolder  
**Core promise:** clone a project on a VPS, run one tool, receive a clean deployment kit you can review, validate, and apply.

depctl is not a PaaS, not a black-box deploy bot, and not a random Dockerfile generator. It is a **deployment-kit compiler**:

> repository + target preset + a few deployment answers → reviewed `.deploy/` kit → safe apply.

## The core workflow

```bash
cd /srv/my-project

depctl scan
depctl plan --preset compose-traefik --domain app.example.com
depctl write
depctl validate
depctl apply
```

Alternative Docker distribution:

```bash
docker run --rm \
  -v "$PWD:/workspace" \
  yourname/depctl:latest scan

docker run --rm \
  -v "$PWD:/workspace" \
  yourname/depctl:latest plan --preset compose-traefik --domain app.example.com
```

`apply` may need host Docker access. That must be explicit and guarded.

## Main docs

Read in this order:

1. [`docs/00_MANIFESTO.md`](docs/00_MANIFESTO.md)
2. [`docs/01_PRODUCT_SPEC.md`](docs/01_PRODUCT_SPEC.md)
3. [`docs/02_STRONG_MVP.md`](docs/02_STRONG_MVP.md)
4. [`docs/03_USER_FLOWS.md`](docs/03_USER_FLOWS.md)
5. [`docs/04_CLI_UX.md`](docs/04_CLI_UX.md)
6. [`docs/05_SCANNER_ENGINE.md`](docs/05_SCANNER_ENGINE.md)
7. [`docs/06_DETECTION_RULES.md`](docs/06_DETECTION_RULES.md)
8. [`docs/07_DEPLOYMENT_KIT_CONTRACT.md`](docs/07_DEPLOYMENT_KIT_CONTRACT.md)
9. [`docs/08_SECURITY_AND_SECRETS.md`](docs/08_SECURITY_AND_SECRETS.md)
10. [`docs/09_APPLY_ROLLBACK_IDEMPOTENCY.md`](docs/09_APPLY_ROLLBACK_IDEMPOTENCY.md)
11. [`docs/10_TEMPLATE_QUALITY_BAR.md`](docs/10_TEMPLATE_QUALITY_BAR.md)
12. [`docs/11_TESTING_STRATEGY.md`](docs/11_TESTING_STRATEGY.md)
13. [`docs/12_ROADMAP.md`](docs/12_ROADMAP.md)
14. [`docs/13_BUILD_CHECKLIST.md`](docs/13_BUILD_CHECKLIST.md)
15. [`docs/15_COMMAND_SET.md`](docs/15_COMMAND_SET.md)
16. [`docs/16_EXPLANATION_REPORTS.md`](docs/16_EXPLANATION_REPORTS.md)

## Project rule

Build one excellent path before adding many presets.

The recommended first path is:

> Single VPS + Docker Compose + Traefik + automatic HTTPS + Node/Laravel/Python detection + `.env.example` + deploy scripts + GitHub Actions template.

Swarm, Portainer, Gitea, monorepos, queues, workers, and advanced secrets should come after the core path is reliable.
