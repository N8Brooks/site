# naterpatater.com

Minimal root-domain site for `naterpatater.com`.

It is intentionally small:

- one static React page for the public site
- AdSense verification snippet in the root HTML
- a tiny Go server for production asset delivery, health checks, and security headers
- Docker and GitHub Actions for CI, image scanning, and GHCR publication

## Local development

Requirements:

- Node 24+
- Go 1.26+

Commands:

```sh
npm install
npm run dev
```

## Verification

```sh
npm run lint
npm run build
go test ./cmd/server ./internal/...
GOCACHE=/tmp/go-build GOMODCACHE=/tmp/gomodcache go build -o /tmp/site ./cmd/server
```

The Go server expects `dist/` to exist, so run `npm run build` before Go build or test commands.

## Container

Build the production image:

```sh
docker build -t site:local .
```

Run it locally:

```sh
docker run --rm -p 8080:8080 site:local
```

## Deployment

GitHub Actions provides:

- `ci`: lint, frontend build, Go tests, and Go build
- `docker`: image build plus Trivy scan
- `release`: multi-arch push to `ghcr.io/<owner>/site` from `main`

## Attribution

See [docs/attribution.md](docs/attribution.md).
