# web

Root-domain site for `naterpatater.com`.

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

Build the production image from the repo root:

```sh
docker build -f web/Dockerfile -t site-web:local .
```

Run it locally:

```sh
docker run --rm -p 8080:8080 site-web:local
```

## Attribution

See [docs/attribution.md](docs/attribution.md).
