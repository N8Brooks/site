# site

Monorepo for `naterpatater.com`.

## Layout

- `web/`: the root-domain site, including the React app, embedded Go server, Dockerfile, and app-specific docs
- `k8s/`: Kubernetes manifests for deploying the site by itself

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
- Canonical homelab deployment repo: `ssh://git@github.com/N8Brooks/homelab-cluster`
- Canonical Flux path: `./clusters/homelab` in the `homelab-cluster` repo
- This repo should only build and publish `site-web`; do not add new homelab-specific deploy state here
