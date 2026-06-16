# Product Specification

## Name

**depctl**

Pronounced like “deploy control.”

The name fits the product because the tool is not a platform, not a dashboard, and not a magic deploy bot. It is a command-line control layer for deployment configuration on a VPS.

## One-liner

depctl scans a cloned project on a VPS and generates a reviewable deployment kit: Docker, reverse proxy, environment, CI/CD, scripts, validation, and safe apply.

## Category

Repo-aware VPS deployment scaffolder.

It sits between:

- Docker Init;
- Nixpacks/Railpack-like builders;
- Ansible;
- Dokku;
- Coolify;
- Kamal;
- Portainer;
- hand-written DevOps.

But its unique angle is:

> generate readable deployment infrastructure from a repo, inside a VPS workflow, without becoming a full platform.

## Primary use case

A developer has a VPS and a project repository.

They want to deploy the project quickly without manually writing all Docker, proxy, SSL, env, and CI files.

## Product thesis

Most VPS deployments need the same decisions:

- What is this app?
- How is it built?
- How is it run?
- Which port does it expose?
- Which services does it need?
- Which env vars does it expect?
- How should traffic reach it?
- How should it be updated?
- How can it be rolled back?

depctl should answer these from repo signals plus a few explicit user choices.

## Target user

### Primary

Developers who own a VPS and deploy small to medium web apps.

Examples:

- solo developer;
- student building projects;
- freelancer hosting client apps;
- small agency;
- technical founder;
- DevOps beginner who still wants real files.

### Secondary

DevOps engineers who want to reduce repetitive setup work.

## Jobs to be done

1. "I cloned a repo on my VPS. Help me deploy it correctly."
2. "Generate Docker and reverse proxy files that match this app."
3. "Create `.env.example` without leaking secrets."
4. "Give me deploy and rollback scripts."
5. "Wire CI/CD for this repo."
6. "Tell me what is missing before I break production."
7. "Let me review everything before applying."

## Product boundaries

depctl should not be:

- a hosting platform;
- a dashboard-first product;
- a Kubernetes platform;
- an AI agent that edits everything;
- a secret manager;
- a monitoring platform;
- a replacement for Docker;
- a replacement for CI providers.

It can integrate with these things later, but it should remain a CLI-first deployment-kit generator.

## Core concepts

### Scan

Read deployment-relevant files and produce `detected.json`.

### Plan

Convert detection plus user choices into `plan.json`.

### Write

Render `.deploy/` files from the plan.

### Validate

Check generated files and host readiness.

### Apply

Run the reviewed plan.

### Rollback

Restore previous working deployment state.

## Non-negotiable product rules

1. Never silently overwrite user files.
2. Never leak secret values into generated files.
3. Never apply a plan that has not been written and validated.
4. Never pretend confidence when detection is weak.
5. Always produce both human-readable and machine-readable output.
6. Always make generated files boring, clean, and editable.
7. Always support dry-run mode.
8. Always make `apply` idempotent.
9. Always create backups before destructive operations.
10. Always explain what the user must do manually when automation is unsafe.

## Product success criteria

The tool is successful when a user says:

> I used this on my own VPS because it saved me an hour and avoided the stupid mistakes I usually make.

Not:

> It generated a lot of files, but I do not understand them.
