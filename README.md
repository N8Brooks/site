# site

Monorepo for `naterpatater.com`.

## Layout

- `web/`: the root-domain site, including the React app, embedded Go server, Dockerfile, and app-specific docs
- `k8s/`: Kubernetes manifests for deploying the site by itself
- `clusters/homelab/`: Flux entrypoint that applies the site manifests from this repo

## Local app development

```sh
cd web
npm install
npm run dev
```

## Verification

```sh
cd web
npm run lint
npm run build
go test ./cmd/server ./internal/...
GOCACHE=/tmp/go-build GOMODCACHE=/tmp/gomodcache go build -o /tmp/site ./cmd/server
```

## Image and deployment

- Docker image: `ghcr.io/<owner>/site-web`
- Flux path for this repo: `./clusters/homelab`
- External Flux repo should add a `GitRepository` for `ssh://git@github.com/N8Brooks/site` and a `Kustomization` that points at `./clusters/homelab`
