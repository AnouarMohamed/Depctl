# Depctl

**Depctl** is a modern deployment controller designed to simplify the lifecycle of containerized applications. It focuses on a **Scan -> Plan -> Apply** workflow to provide predictability and safety in deployments.

## Features

- **Automated Scanning:** Detects project types (Node.js, Go, Python, PHP, etc.), databases, and environment requirements.
- **Predictive Planning:** Generates a deployment plan and a detailed report before any changes are made.
- **Idempotent Application:** Safely applies configurations and manages rollbacks.
- **Extensible Architecture:** Easily add new language detectors and deployment targets.

## Quick Start

### 1. Installation

#### Binary
```bash
go install github.com/AnouarMohamed/Depctl@latest
```

#### Docker
You can also run `depctl` as a Docker container without installing Go:

```bash
docker pull anouarmohamedx/depctl:latest

# Run depctl (e.g., scan the current directory)
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):/workspace \
  anouarmohamedx/depctl:latest scan /workspace
```

### 2. Scan your project
...

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
