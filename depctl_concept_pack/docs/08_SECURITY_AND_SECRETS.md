# Security and Secrets

## Security principle

depctl should reduce deployment mistakes, not become a new risk.

## Default security posture

- Read-only scan by default.
- Write only inside `.deploy/`.
- No secret values copied from `.env`.
- No Docker socket required for scan/plan/write.
- Docker socket required only for apply when running inside a container.
- `apply` requires confirmation.
- Destructive operations require explicit flags.
- Backups before overwrites.
- Database volumes are never deleted by default.

## Docker image distribution risk

The user wants to run depctl as a Docker image from Docker Hub.

That is convenient.

But for apply mode, mounting Docker socket is powerful:

```bash
-v /var/run/docker.sock:/var/run/docker.sock
```

This effectively gives the container high control over the host.

Therefore:

- `scan`, `plan`, `write`, `validate` should not require Docker socket.
- `apply` should clearly warn when Docker socket is mounted.
- host binary should be recommended for apply if possible.
- Docker image should run as non-root where possible.

## Secret handling

### What to do

- collect env var names;
- generate `.env.example` with empty values;
- flag sensitive names;
- document required secrets;
- optionally generate Docker secrets templates later.

### What not to do

- never copy values from `.env`;
- never print secret values;
- never write real secrets to compose files;
- never commit generated `.env`;
- never infer fake secret values.

## `.env.example` format

Good:

```env
APP_ENV=production
APP_URL=
DATABASE_URL=
REDIS_URL=
APP_SECRET=
```

Bad:

```env
APP_SECRET=change-me-super-secret
DB_PASSWORD=password123
```

Use blank values or safe placeholders only when necessary.

## Sensitive key detection

Flag keys containing:

```text
SECRET
TOKEN
PASSWORD
PASS
PRIVATE
KEY
API_KEY
DB_PASS
JWT
SESSION
CREDENTIAL
ACCESS_KEY
```

## Compose security checks

Warn if:

- database port exposed to host;
- Redis port exposed to host;
- app runs as root;
- no restart policy;
- secrets are hardcoded;
- `.env` is referenced but missing;
- Docker socket is mounted into app service;
- privileged mode is used;
- host filesystem mounted broadly.

## Reverse proxy security

For Traefik:

- expose only proxy ports 80/443;
- route app through internal Docker network;
- do not expose database;
- generate HTTPS labels;
- use secure entrypoints.

For Nginx:

- no direct database exposure;
- redirect HTTP to HTTPS when configured;
- document Certbot/SSL step.

## CI/CD security

Generated CI should:

- require SSH host, user, key, and path as provider secrets;
- not include private keys directly;
- use known host verification where possible;
- avoid `curl | bash`;
- run a fixed deploy script on server;
- fail fast.

## Dangerous operations

These should require explicit confirmation:

- deleting containers;
- deleting images;
- deleting volumes;
- changing firewall rules;
- overwriting existing reverse proxy config;
- running migrations;
- opening ports;
- mounting Docker socket.

## Security MVP checklist

- [ ] No secret values copied.
- [ ] `.env.example` generated safely.
- [ ] Docker socket not used for scan/write.
- [ ] Apply confirmation required.
- [ ] Existing files backed up.
- [ ] Database ports not exposed by default.
- [ ] Generated CI uses secrets.
- [ ] Validation flags common security mistakes.
