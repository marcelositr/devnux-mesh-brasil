#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

cp "$repo_root/.dockerignore" "$tmp/.dockerignore"
printf 'sentinel\n' > "$tmp/.env"
printf 'sentinel\n' > "$tmp/.env.production"
printf 'example configuration\n' > "$tmp/.env.example"
mkdir -p "$tmp/certs" "$tmp/secrets"
printf 'sentinel\n' > "$tmp/certs/private.key"
printf 'sentinel\n' > "$tmp/secrets/token"

docker build --no-cache -f - "$tmp" <<'DOCKERFILE'
FROM alpine:3.22
COPY . /context
RUN test ! -e /context/.env \
    && test ! -e /context/.env.production \
    && test ! -e /context/certs/private.key \
    && test ! -e /context/secrets/token \
    && test -f /context/.env.example
DOCKERFILE
