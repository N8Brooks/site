FROM node:24-alpine AS frontend-builder

WORKDIR /workspace/web

COPY web/package.json web/package-lock.json ./

RUN npm install

COPY web/index.html ./
COPY web/tsconfig.json web/tsconfig.app.json web/tsconfig.node.json ./
COPY web/vite.config.ts ./
COPY web/src ./src
COPY web/public ./public

RUN npm run build

FROM golang:1.26.2 AS server-builder

WORKDIR /workspace

COPY web ./web
COPY --from=frontend-builder /workspace/web/dist ./web/dist

RUN cd web && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/site ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=server-builder /out/site /usr/local/bin/site

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/site"]
