.PHONY: test coverage coverage-html coverage-profile

test:
	go test ./...

coverage-profile:
	@tmp="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp"' EXIT; \
	merge_set() { \
		out="$$1"; \
		shift; \
		{ \
			echo 'mode: set'; \
			awk 'FNR==1 { next } { key=$$1 " " $$2; if ($$3+0 > seen[key]) seen[key]=$$3+0 } END { for (key in seen) print key, (seen[key] > 0 ? 1 : 0) }' "$$@" | sort; \
		} > "$$out"; \
	}; \
	go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess -coverprofile="$$tmp/root.cover" >/dev/null; \
	go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/game -coverprofile="$$tmp/game.cover" >/dev/null; \
	go test ./internal/domain/chess/internal/bitboard -coverpkg=./internal/domain/chess/internal/bitboard -coverprofile="$$tmp/bitboard.cover" >/dev/null; \
	go test ./internal/domain/chess/internal/geom -coverpkg=./internal/domain/chess/internal/geom -coverprofile="$$tmp/geom.cover" >/dev/null; \
	go test ./internal/domain/chess/model -coverpkg=./internal/domain/chess/model -coverprofile="$$tmp/model-int.cover" >/dev/null; \
	go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/model -coverprofile="$$tmp/model-ext.cover" >/dev/null; \
	go test ./internal/domain/chess/position -coverpkg=./internal/domain/chess/position -coverprofile="$$tmp/position-int.cover" >/dev/null; \
	go test ./internal/tests/domain/chess -coverpkg=./internal/domain/chess/position -coverprofile="$$tmp/position-ext.cover" >/dev/null; \
	merge_set "$$tmp/model.cover" "$$tmp/model-int.cover" "$$tmp/model-ext.cover"; \
	merge_set "$$tmp/position.cover" "$$tmp/position-int.cover" "$$tmp/position-ext.cover"; \
	{ \
		echo 'mode: set'; \
		tail -n +2 "$$tmp/root.cover"; \
		tail -n +2 "$$tmp/game.cover"; \
		tail -n +2 "$$tmp/bitboard.cover"; \
		tail -n +2 "$$tmp/geom.cover"; \
		tail -n +2 "$$tmp/model.cover"; \
		tail -n +2 "$$tmp/position.cover"; \
	} > coverage.out

coverage: coverage-profile
	@report="$$(go tool cover -func=coverage.out)"; \
	echo "$$report"; \
	echo "$$report" | grep -Eq '^total:.*100\.0%$$'

coverage-html: coverage-profile
	go tool cover -html=coverage.out -o coverage.html
