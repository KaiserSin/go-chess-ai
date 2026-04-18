#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

cd "$ROOT_DIR"

go test -tags=extended ./internal/application/ai ./internal/tests/application/ai
