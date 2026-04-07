# syntax=docker/dockerfile:1.5
# =============================================================================
# Startr/feeds — two-stage Go build
#   1. go-builder:  Go 1.25 — compile ./cmd/feeds from local source
#   2. runtime:     Alpine 3.21 — just the binary (~15 MB)
#
# Adapted from WEB-DB-sage-pb. The v0.2 release will swap the long-running
# `feeds serve` subcommand's stdlib ticker for a PocketBase framework import,
# which is why the runtime stage preserves the pb_data volume and 8090 port
# contract. v0.1.0 serve mode ignores them but CapRover deployments stay
# forward-compatible.
# =============================================================================

# =============================================================================
# Stage 1: GO — compile the feeds binary from local source
# =============================================================================
FROM golang:1.25-alpine AS go-builder

ARG TARGETARCH

WORKDIR /src

# Copy module files first (cacheable layer). go.sum is generated on first
# build via `go mod download` if missing — the [m] glob allows it to be absent.
COPY go.mod go.su[m] ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Run the 5 critical tests gating v0.1.0. If any fail, the build fails loud
# and nothing ships. `go vet` is cheap and catches a lot of silly mistakes
# so it runs too.
RUN go vet ./... && go test ./...

# Build the feeds binary (statically linked, no CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w -X main.version=v0.1.0" -o /feeds ./cmd/feeds

# =============================================================================
# Stage 2: RUNTIME — minimal Alpine with just the binary
# =============================================================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=go-builder /feeds /usr/local/bin/feeds

WORKDIR /app

# Data volume — used by `feeds serve` in v0.2+ (PocketBase).
# v0.1.0 serve mode writes state to /app/pb_data/.feeds-state.json.
VOLUME /app/pb_data

EXPOSE 8090

CMD ["feeds", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
