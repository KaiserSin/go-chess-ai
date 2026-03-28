#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

cd "$ROOT_DIR"

go vet ./...

if ! command -v golangci-lint >/dev/null 2>&1; then
	echo "golangci-lint is required but not installed. Install it and re-run ./scripts/lint.sh." >&2
	exit 127
fi

exec golangci-lint run ./...
