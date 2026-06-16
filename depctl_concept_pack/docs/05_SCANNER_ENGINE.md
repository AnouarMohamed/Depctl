# Scanner Engine

## Purpose

The scanner converts a project directory into deployment-relevant facts.

It should not understand the whole codebase.

It should only find what is needed to deploy safely.

## Scanner outputs

The scanner writes:

```text
.deploy/detected.json
.deploy/reports/scan-report.md
```

## Detection model

Every detected fact should include:

- value;
- source file;
- confidence;
- reason;
- optional warnings.

Example:

```json
{
  "runtime": {
    "name": "node",
    "confidence": 0.94,
    "evidence": [
      {
        "file": "package.json",
        "reason": "package.json exists with scripts.start"
      }
    ]
  }
}
```

## Confidence levels

Use simple bands:

```text
90-100: strong
70-89: likely
50-69: weak
0-49: unknown
```

The tool should not generate risky config from weak detection without asking.

## Scanner layers

### Layer 1 — repository structure

Detect:

- single app;
- monorepo;
- multiple service folders;
- packages/apps/services/src structure;
- git root;
- existing `.deploy/`.

### Layer 2 — runtime and framework

Detect:

- Node/Bun;
- Laravel/PHP;
- Python;
- Go;
- Ruby/Java later.

### Layer 3 — build and run commands

Detect:

- build command;
- start command;
- dev command;
- package manager;
- compiled output path;
- public/static path.

### Layer 4 — network

Detect:

- app port;
- exposed port in Dockerfile;
- env port usage;
- framework default port;
- existing reverse proxy config.

### Layer 5 — dependencies

Detect:

- PostgreSQL;
- MySQL/MariaDB;
- Redis;
- MongoDB later;
- queues;
- scheduler/cron;
- file uploads/persistent storage.

### Layer 6 — containerization

Detect:

- Dockerfile;
- Dockerfile quality;
- Compose file;
- compose services;
- volumes;
- networks;
- healthchecks;
- exposed ports.

### Layer 7 — env and secrets

Detect:

- `.env`;
- `.env.example`;
- env var names in code;
- risky secret names;
- missing example keys.

Never copy secret values.

### Layer 8 — CI/CD

Detect:

- `.github/workflows`;
- `.gitlab-ci.yml`;
- `.gitea/workflows`;
- Jenkinsfile;
- deployment scripts.

## Scanner architecture

Suggested Go packages:

```text
internal/scanner/
  scanner.go
  evidence.go
  confidence.go
  runtime_node.go
  runtime_php.go
  runtime_python.go
  runtime_go.go
  dockerfile.go
  compose.go
  env.go
  ci.go
```

## Scanner rule format

Each rule should be small and testable.

Example pseudo-structure:

```go
type Detection struct {
    Key        string
    Value      any
    Confidence float64
    Evidence   []Evidence
    Warnings   []Warning
}
```

## Important scanner rule

Do not let one file decide everything.

Example:

- `package.json` means Node project.
- `next` dependency means likely Next.js.
- `scripts.start` confirms runtime command.
- `next.config.js` increases confidence.

## Scanner should produce warnings

Examples:

- no start command;
- no exposed port;
- env vars found but no `.env.example`;
- existing Dockerfile has no non-root user;
- database dependency found but no database config;
- multiple apps detected but no public service selected.

## Do not scan deeply in MVP

Avoid parsing whole codebase if simple file-level detection is enough.

MVP can inspect common config files and limited source patterns.
