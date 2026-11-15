# Kaqen

**One logo. Every asset.**

Open-source, self-hostable logo asset generator by Ploos AS.

## Status
Bootstrap only. Upload and asset generation are not yet implemented.

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
