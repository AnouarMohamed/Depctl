# depctl Manifesto

## 1. The problem

Deploying a normal web app on a VPS is still too manual.

A developer clones a repo, then repeats the same fragile work:

- identify the framework;
- guess the port;
- write or fix a Dockerfile;
- create a compose file;
- configure Nginx or Traefik;
- handle HTTPS;
- create `.env.example`;
- avoid leaking secrets;
- write deploy and rollback scripts;
- wire CI/CD;
- test everything;
- hope nothing was missed.

This work is not intellectually hard every time. It is just easy to get wrong.

depctl exists to remove the repetitive mistakes from VPS deployment work while keeping the result understandable.

## 2. The belief

Good DevOps automation should not hide everything.

It should produce clean, readable, boring infrastructure files that a human can review.

The user should never wonder:

> What did this tool do to my server?

Instead, the tool should say:

> Here is exactly what I detected.  
> Here is exactly what I plan to generate.  
> Here are the files.  
> Review them.  
> Validate them.  
> Then apply them.

## 3. The product

depctl is a **repo-aware deployment-kit compiler for VPS projects**.

It reads only deployment-relevant signals from a project directory, then generates a `.deploy/` kit containing:

- Dockerfile or Dockerfile patch suggestions;
- Docker Compose or stack manifests;
- Traefik/Nginx reverse proxy config;
- `.env.example`;
- deploy, rollback, status, and backup scripts;
- CI/CD workflow templates;
- a machine-readable deployment plan;
- a human-readable report.

It does not pretend to be magic.

It turns messy VPS deployment work into a repeatable, reviewable, documented process.

## 4. The safety principle

The default mode is read-only.

The generation mode writes only to `.deploy/`.

The apply mode runs only from a reviewed plan.

`apply` must not rescan and improvise. It must apply the plan the user reviewed.

## 5. The quality principle

Generated files must look like a good DevOps engineer wrote them manually.

No tutorial-grade Dockerfiles.  
No random options.  
No half-configured services.  
No mysterious generated clutter.  
No overwrites without backup.  
No fake confidence.

If the tool is not sure, it must say so.

## 6. The ideal user experience

A developer should be able to SSH into a VPS and do this:

```bash
cd /srv/app
depctl plan --preset compose-traefik --domain app.example.com
depctl write
depctl validate
depctl apply
```

Then get a working deployment with readable files and no guessing.

## 7. The promise

depctl does not replace DevOps knowledge.

It packages the boring, repeated, error-prone parts of DevOps into a tool that works the same way every time.

Its job is not to be clever.

Its job is to be correct.
