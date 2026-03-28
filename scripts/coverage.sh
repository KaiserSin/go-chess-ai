#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
HTML_OUTPUT=0

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

merge_set() {
	out_file=$1
	shift

	{
		echo 'mode: set'
		awk 'FNR==1 { next } { key=$1 " " $2; if ($3+0 > seen[key]) seen[key]=$3+0 } END { for (key in seen) print key, (seen[key] > 0 ? 1 : 0) }' "$@" | sort
	} > "$out_file"
}

go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess -coverprofile="$tmp_dir/root.cover" >/dev/null
go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/game -coverprofile="$tmp_dir/game.cover" >/dev/null
go test ./internal/domain/chess/internal/bitboard -coverpkg=./internal/domain/chess/internal/bitboard -coverprofile="$tmp_dir/bitboard.cover" >/dev/null
go test ./internal/domain/chess/internal/geom -coverpkg=./internal/domain/chess/internal/geom -coverprofile="$tmp_dir/geom.cover" >/dev/null
go test ./internal/domain/chess/model -coverpkg=./internal/domain/chess/model -coverprofile="$tmp_dir/model-int.cover" >/dev/null
go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/model -coverprofile="$tmp_dir/model-ext.cover" >/dev/null
go test ./internal/domain/chess/position -coverpkg=./internal/domain/chess/position -coverprofile="$tmp_dir/position-int.cover" >/dev/null
go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/position -coverprofile="$tmp_dir/position-ext.cover" >/dev/null

merge_set "$tmp_dir/model.cover" "$tmp_dir/model-int.cover" "$tmp_dir/model-ext.cover"
merge_set "$tmp_dir/position.cover" "$tmp_dir/position-int.cover" "$tmp_dir/position-ext.cover"

{
	echo 'mode: set'
	tail -n +2 "$tmp_dir/root.cover"
	tail -n +2 "$tmp_dir/game.cover"
	tail -n +2 "$tmp_dir/bitboard.cover"
	tail -n +2 "$tmp_dir/geom.cover"
	tail -n +2 "$tmp_dir/model.cover"
	tail -n +2 "$tmp_dir/position.cover"
} > coverage.out

report=$(go tool cover -func=coverage.out)
printf '%s\n' "$report"
printf '%s\n' "$report" | grep -Eq '^total:.*100\.0%$'

if [ "$HTML_OUTPUT" -eq 1 ]; then
	go tool cover -html=coverage.out -o coverage.html
fi
