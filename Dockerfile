FROM node:24-alpine AS frontend-builder

WORKDIR /workspace

COPY package.json package-lock.json ./

RUN npm install

COPY index.html ./
COPY tsconfig.json tsconfig.app.json tsconfig.node.json ./
COPY vite.config.ts ./
COPY src ./src
COPY public ./public

RUN npm run build

FROM golang:1.26.2 AS server-builder

WORKDIR /workspace

COPY go.mod ./
COPY assets.go ./
COPY cmd ./cmd
COPY internal ./internal
COPY --from=frontend-builder /workspace/dist ./dist

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/site ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=server-builder /out/site /usr/local/bin/site

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/site"]
