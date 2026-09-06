FROM golang:1.27.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=.git,target=/app/.git,ro \
    CGO_ENABLED=1 GOOS=linux go build \
    -tags "fts5" \
    -ldflags="-X bananajeanss/go-ship/handlers.commitHash=$(git rev-parse --short HEAD)" \
    -o /app/go-ship

FROM alpine:3.24.1
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app /app

RUN mkdir -p /data

CMD ["./go-ship"]