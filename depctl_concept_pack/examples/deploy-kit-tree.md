# Example `.deploy/` Tree

```text
.deploy/
  README.md
  detected.json
  plan.json
  Dockerfile
  .dockerignore
  docker-compose.yml
  .env.example
  ci/
    github-actions.yml
  scripts/
    deploy.sh
    rollback.sh
    status.sh
    backup.sh
  reports/
    scan-report.md
    plan-report.md
    validation-report.md
  backups/
```

## Rule

This folder is the contract between scanner, planner, writer, validator, and apply engine.

Do not let `apply` use hidden state.
