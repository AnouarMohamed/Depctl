# Depctl

**Depctl** is a modern deployment controller designed to simplify the lifecycle of containerized applications. It focuses on a **Scan -> Plan -> Apply** workflow to provide predictability and safety in deployments.

## Features

- **Automated Scanning:** Detects project types (Node.js, Go, Python, PHP, etc.), databases, and environment requirements.
- **Predictive Planning:** Generates a deployment plan and a detailed report before any changes are made.
- **Idempotent Application:** Safely applies configurations and manages rollbacks.
- **Extensible Architecture:** Easily add new language detectors and deployment targets.

## Quick Start

### 1. Installation

```bash
go install github.com/AnouarMohamed/Depctl@latest
```

### 2. Scan your project

Run the scanner to detect your application structure:

```bash
depctl scan
```

### 3. Generate a plan

Review the proposed deployment strategy:

```bash
depctl plan
```

### 4. Apply the deployment

Execute the plan once you are satisfied with the report:

```bash
depctl apply
```

## Project Structure

- `cmd/`: CLI command implementations.
- `internal/scanner/`: Core detection logic for various languages and environments.
- `internal/planner/`: Logic for generating deployment plans and reports.
- `depctl_concept_pack/`: Detailed documentation and architectural decision logs.

## Documentation

For deep dives into the architecture and design principles, see the [Concept Pack](./depctl_concept_pack/README.md).
