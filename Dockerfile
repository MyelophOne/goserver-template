FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS="-trimpath" go build -pgo=auto -ldflags="-s -w -X github.com/myelophone/goserver.AppEnv=prod -X github.com/myelophone/goserver.AppVersion=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" -o /app/goserver ./cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS="-trimpath" go build -ldflags="-s -w" -o /app/healthcheck ./cmd/healthcheck

FROM gcr.io/distroless/static:nonroot

WORKDIR /home/nonroot/app

COPY --from=builder /app/assets ./assets
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/goserver .
COPY --from=builder /app/healthcheck .
COPY --from=builder /app/LICENSE ./LICENSE

LABEL org.opencontainers.image.title="@myelophone/goserver-template"
LABEL org.opencontainers.image.description="High-performance Go server app based on @myelophone/goserver by @myeloph.one"
LABEL org.opencontainers.image.authors="Aliaksandr Ivanou"
LABEL org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0"
LABEL org.opencontainers.image.vendor="Aliaksandr Ivanou"
LABEL org.opencontainers.image.source="https://github.com/myelophone/goserver-template"

EXPOSE 8080

CMD ["./goserver"]
