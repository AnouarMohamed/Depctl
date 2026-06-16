# Explanation and Temporary Reports

## Why this matters

depctl should not only generate deployment files. It should explain the deployment decisions it made.

A user should be able to answer:

- What did depctl detect?
- Why did it choose this runtime?
- Why did it expose this port?
- Why did it generate Traefik labels or Nginx config?
- Which env vars are missing?
- Which secrets must I fill manually?
- What exactly will `apply` do?
- What could break?

## Report-first design

Every major command should produce two outputs:

1. machine-readable JSON for the tool;
2. human-readable Markdown for the user.

```text
.deploy/
  detected.json
  plan.json
  reports/
    scan-report.md
    plan-report.md
    validation-report.md
    apply-report.md
```

## Temporary report mode

Sometimes the user wants explanations without keeping files in the repo.

Support:

```bash
depctl scan --explain
```

This prints a summary to the terminal and writes a temporary Markdown report:

```text
/tmp/depctl/<project-name>/scan-report.md
```

Support:

```bash
depctl setup --reports-only
```

This generates reports but does not generate deployment files.

## Persistent report mode

Default project mode writes reports into `.deploy/reports/` because those reports are part of the deployment kit.

```bash
depctl scan
depctl plan --preset compose-traefik --domain app.example.com
depctl validate
```

Expected result:

```text
.deploy/reports/scan-report.md
.deploy/reports/plan-report.md
.deploy/reports/validation-report.md
```

## Suggested report contents

### `scan-report.md`

Should include:

- project path;
- detected runtime and framework;
- confidence score;
- evidence used;
- detected commands;
- detected ports;
- detected services;
- detected env vars;
- warnings;
- unsupported or uncertain items.

Example:

```md
# Scan Report

Detected: Node / Next.js
Confidence: 91%

Evidence:
- package.json contains `next`
- scripts contain `build` and `start`
- default Next.js port 3000 used because no explicit PORT was found

Warnings:
- DATABASE_URL is referenced but no .env.example exists
- no health endpoint detected
```

### `plan-report.md`

Should include:

- selected preset;
- domain;
- generated services;
- reverse proxy choice;
- volumes;
- networks;
- generated files;
- manual secrets required;
- exact next command.

### `validation-report.md`

Should include:

- passed checks;
- warnings;
- blocking errors;
- recommended fixes;
- whether `apply` is allowed.

### `apply-report.md`

Should include:

- timestamp;
- plan applied;
- commands executed;
- containers started;
- healthcheck result;
- rollback snapshot location.

## Terminal explanation

Terminal output should be short. Markdown reports should be detailed.

Good terminal output:

```text
Scan complete: Node / Next.js detected with 91% confidence.
Report written: .deploy/reports/scan-report.md
Next: depctl plan --preset compose-traefik --domain app.example.com
```

Bad terminal output:

```text
Huge wall of logs and generated YAML explanation.
```

## `depctl explain`

A later command can expose the explanation system directly:

```bash
depctl explain scan
depctl explain plan
depctl explain warnings
depctl explain files
```

This command should read existing `.deploy/*.json` and reports. It should not rescan unless explicitly requested.

## The rule

Explanations must be evidence-based.

If depctl does not know something, it should say so clearly and ask for a flag or prompt.
