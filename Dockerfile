FROM node:26-alpine AS web-deps

FROM golang:alpine AS builder

RUN apk add --no-cache git libstdc++

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY --from=web-deps /usr/local/bin/node /usr/local/bin/node

RUN APP_ENV=prod GIT_COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
    go run -tags webcli ./cmd generate \
    && APP_ENV=prod GIT_COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
    go run -tags webcli ./cmd build

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates \
    && addgroup -S goserver \
    && adduser -S -G goserver -h /app goserver

COPY --chown=goserver:goserver --from=builder /app/dist/ ./
COPY --chown=goserver:goserver --from=builder /app/LICENSE ./LICENSE

USER goserver

LABEL org.opencontainers.image.title="@myelophone/goserver-template"
LABEL org.opencontainers.image.description="Production-ready Go application template with optional SSR web mode by @myeloph.one"
LABEL org.opencontainers.image.authors="Aliaksandr Ivanou"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"
LABEL org.opencontainers.image.vendor="Aliaksandr Ivanou"
LABEL org.opencontainers.image.source="https://github.com/myelophone/goserver-template"

EXPOSE 8080

CMD ["./goserver"]
