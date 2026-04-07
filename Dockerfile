# =============================================================================
# Startr/feeds — PocketBase standalone + JS hooks
#
# No Go compiler. No dependency resolution. No compilation.
# Just download the pre-built PocketBase binary and copy JS hooks.
# Build time: seconds.
#
# The rewrite pipeline runs as a PocketBase JS hook (pb_hooks/feeds.pb.js).
# PocketBase provides: HTTP server, static file serving (pb_public/),
# admin UI at /_/, cron scheduler, SQLite, and graceful shutdown.
#
# The Go pipeline code in internal/ is kept as a reference implementation
# and high-performance fallback for large-scale multi-feed deployments.
# =============================================================================

FROM alpine:3.21

ARG PB_VERSION=0.36.8
ARG TARGETARCH

RUN apk add --no-cache ca-certificates tzdata unzip

# Download pre-built PocketBase binary and rename to "feeds" for branding.
ADD https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_${TARGETARCH}.zip /tmp/pb.zip
RUN unzip /tmp/pb.zip -d /tmp/ \
    && mv /tmp/pocketbase /usr/local/bin/feeds \
    && chmod +x /usr/local/bin/feeds \
    && rm -f /tmp/pb.zip /tmp/CHANGELOG.md /tmp/LICENSE.md

WORKDIR /app

# JS hooks — the rewrite pipeline lives here.
# Libraries in pb_hooks/lib/ are vendored dist bundles (no npm install).
# See pb_hooks/lib/LIBRARIES.md for how to add/update libraries.
COPY pb_hooks/ ./pb_hooks/

# Migrations — creates the "feeds" collection on first boot.
COPY pb_migrations/ ./pb_migrations/

# Static assets served by PocketBase. The rewrite pipeline writes output
# XML here (e.g., pb_public/v1/show.xml → https://feed.example.com/v1/show.xml).
COPY pb_public/ ./pb_public/

# Data volume — PocketBase SQLite database + feeds cache state.
VOLUME /app/pb_data

EXPOSE 8090

# PB defaults to --http=0.0.0.0:8080; we override to 8090 to match
# existing CapRover / Makefile infrastructure. Feed config comes from
# the "feeds" collection in the admin UI, or FEEDS_* env vars as fallback.
CMD ["feeds", "serve", "--http=0.0.0.0:8090"]
