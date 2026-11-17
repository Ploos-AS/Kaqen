# Kaqen

**One logo. Every asset.**

Open-source, self-hostable logo asset generator by Ploos AS.

## Status
Working Phase A implementation: upload PNG, WebP or restricted static SVG logos, preview icons, and export a ZIP with favicons, PWA/maskable icons, an Open Graph image, HTML integration tags and a setup guide. Pre-public-beta: security and operational verification are still required.

## Local development
Requires Go 1.24+.

```sh
go run ./cmd/kaqen
curl http://localhost:8080/healthz
go test ./...
```

## Docker
```sh
docker compose up --build
```

Compose binds the service to localhost. Put Caddy in front for public HTTPS deployment.

## Roadmap
See [ROADMAP.md](ROADMAP.md). Phase A generates web/PWA assets from an existing logo; Phase B adds a broader brand toolkit.

## License
MIT, see [LICENSE](LICENSE).

See [public beta acceptance checklist](docs/BETA-ACCEPTANCE.md) for manual browser, accessibility, privacy and operational verification.
\n## Beta deployment checklist

- Run `go test ./...`, `go vet ./...`, and `go build ./cmd/kaqen`.
- Put the service behind HTTPS with a reverse proxy; do not expose the development listener directly.
- Configure reverse-proxy request-size limits and request timeouts, and monitor CPU/memory under concurrent uploads.
- Verify exported ZIP contents using PNG, WebP, and static SVG test logos.
- Check that generated `head.html` uses real line breaks and replace the example.com URL.
- Check that `README.txt` has real line breaks and all manifest names are customized.
- Do not log uploaded file bodies; verify actual privacy behavior before public claims.
- Future production deployments should use Ploos-Deploy once it is ready; this repository does not require it yet.

## Runtime safeguards

- HTTP server: 5s header, 30s read, 90s write and 60s idle timeout; 1 MiB maximum request headers.
- At most two image requests run concurrently per process; additional requests receive HTTP 503 and `Retry-After: 5`.
- `GET /healthz` checks that the process is responding; `GET /readyz` confirms the HTTP handler is available. Neither probes downstream infrastructure or actual conversion capacity.
- Reverse proxy should enforce compatible body size and timeout limits. Load-test under real VPS constraints before public beta.
- Ploos-Deploy integration is intentionally deferred.
