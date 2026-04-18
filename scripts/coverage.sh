#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
HTML_OUTPUT=0
MIN_COVERAGE=75.0

case "${1-}" in
	"")
		;;
	html|--html)
		HTML_OUTPUT=1
		;;
	*)
		echo "usage: $0 [html|--html]" >&2
		exit 2
		;;
esac

cd "$ROOT_DIR"

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT HUP INT TERM

run_package() {
	name=$1
	coverpkg=$2
	profile=$3
	shift 3

	go test "$@" -coverpkg="$coverpkg" -coverprofile="$profile" >/dev/null

	report=$(go tool cover -func="$profile")
	total=$(printf '%s\n' "$report" | awk '/^total:/ { print $NF }')
	printf '%-10s %s\n' "$name" "$total"
}

printf 'Package coverage:\n'
run_package domain ./internal/domain/chess "$tmp_dir/domain.cover" ./internal/tests/domain/chess
run_package gameplay ./internal/application/gameplay "$tmp_dir/gameplay.cover" ./internal/tests/application/gameplay ./internal/application/gameplay
run_package ai ./internal/application/ai "$tmp_dir/ai.cover" ./internal/tests/application/ai ./internal/application/ai

{
	echo 'mode: set'
	tail -n +2 "$tmp_dir/domain.cover"
	tail -n +2 "$tmp_dir/gameplay.cover"
	tail -n +2 "$tmp_dir/ai.cover"
} > coverage.out

printf '\nCombined coverage:\n'
report=$(go tool cover -func=coverage.out)
printf '%s\n' "$report" | tail -n 1

total_coverage=$(printf '%s\n' "$report" | awk '/^total:/ { gsub("%", "", $NF); print $NF }')
if [ -z "$total_coverage" ]; then
	echo "failed to parse total coverage from go tool cover output" >&2
	exit 1
fi

if ! awk -v actual="$total_coverage" -v minimum="$MIN_COVERAGE" 'BEGIN { exit ((actual + 0) >= (minimum + 0) ? 0 : 1) }'; then
	echo "coverage check failed: total coverage ${total_coverage}% is below ${MIN_COVERAGE}%" >&2
	exit 1
fi

if [ "$HTML_OUTPUT" -eq 1 ]; then
	go tool cover -html=coverage.out -o coverage.html
fi
