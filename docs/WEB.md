# Kaqen web interface

The Go HTTP server embeds and serves a responsive HTML/CSS/JavaScript interface at `/`. It calls the existing `POST /api/v1/favicons` endpoint. No separate frontend server, build tooling or persistent user storage is needed.

## Local development

```sh
go run ./cmd/kaqen
# Open http://localhost:8080
```

The browser UI accepts PNG and WebP (up to 5 MiB and 16 million pixels), previews the image and downloads a favicon ZIP. SVG validation exists in the inspection API but **SVG rasterization is not yet available**.

## Deployment

Build and publish a versioned Docker image. For now deploy with Docker Compose and Caddy. A future integration with Ploos-Deploy should use a pinned GHCR digest, health check and explicit rollback; do not depend on Ploos-Deploy at runtime.

## Follow-ups

Add end-to-end browser tests, accessible keyboard drag-and-drop verification, image/background/padding controls, maskable PWA icons, SVG rasterization, and upload stress tests before public launch.
