# Forseti deployment readiness (not a deployment runbook)

Target host: `forseti.ousdal.net` (aliases `apps.ousdal.net`, `apps.ploos.no`).
Target app hostname: `beta.kaqen.com`. Production `kaqen.com` and `www.kaqen.com` stay inactive until launch approval.

## Current compatibility boundary

Ploos-Deploy's `docs/architecture.md` is a **target design**, not an implemented production contract. Its current CLI is a disposable Balder test harness. It explicitly says **never adopt unmanaged Compose projects**. Therefore, do **not** claim this Compose project can be adopted automatically, and do not run Ploos-Deploy against Forseti yet.

This repository intentionally avoids fixed container names, shared host ports, hard-coded hostnames and embedded secrets. The Compose project name, host port and image reference are configurable so a future supported Ploos-Deploy migration can reuse the host and app identity without relocating Kaqen.

## Preflight on Forseti (manual, before deployment)

1. Confirm Debian/Docker/Compose/Caddy versions, disk/RAM, firewall, current port allocations, backup and rollback approach.
2. Confirm A and AAAA records for `beta.kaqen.com` resolve to Forseti; check external IPv4 **and** IPv6 HTTPS reachability.
3. Reserve a unique Compose project name (`kaqen-beta`) and an unused loopback port (example `18084`; check before use).
4. Build and publish a versioned Kaqen image to GHCR with a pinned digest. The default `kaqen:local` is for local development only.
5. Validate `docker compose config` and review the rendered output. Do not include secrets in tracked files.
6. Validate Caddy configuration before routing `beta.kaqen.com` to `127.0.0.1:<reserved port>`.
7. Arrange monitoring of `/healthz` and `/readyz`, and verify the service's CPU/memory limits under actual logo conversion load.
8. Complete `docs/BETA-ACCEPTANCE.md` before exposing the beta.

## Compose contract

- Set `COMPOSE_PROJECT_NAME=kaqen-beta` (Compose native variable).
- Set `KAQEN_BIND_PORT` to the verified, unused loopback port.
- Set `KAQEN_IMAGE` to the exact published image digest, e.g. `ghcr.io/ploos-as/kaqen@sha256:<verified-digest>`.
- Use `docker compose --no-ansi config` to review resolved settings; for a prebuilt image use `docker compose up -d --no-build` only after the operational preflight and explicit authorization.
- The app port is bound to IPv4 loopback. Caddy should proxy to that loopback listener, including when clients connect over IPv6.
- `/readyz` is a process readiness endpoint, not a guarantee that image conversion capacity is available.
- Do not mount the Docker socket, publish port 8080 publicly, or run the container privileged.

## Future Ploos-Deploy migration

Implement an explicit, audited **import/adoption workflow** in Ploos-Deploy (or a controlled one-time migration) before attempting takeover. It must verify Compose project identity, images, current release, port ownership, Caddy route, state directory and health, and preserve a tested rollback path. Merely matching directory names is not sufficient. Do not allow Ploos-Deploy to overwrite an unmanaged project.

This document does not create directories or change DNS, Caddy or Docker on Forseti.
