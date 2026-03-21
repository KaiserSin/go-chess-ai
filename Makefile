.PHONY: test coverage coverage-html

DOMAIN_TEST_PKGS := ./internal/domain/chess ./internal/tests/domain/chess
DOMAIN_COVER_PKG := ./internal/domain/chess

test:
	go test ./...

coverage:
	go test $(DOMAIN_TEST_PKGS) -coverpkg=$(DOMAIN_COVER_PKG) -coverprofile=coverage.out
	@report="$$(go tool cover -func=coverage.out)"; \
	echo "$$report"; \
	echo "$$report" | grep -Eq '^total:.*100\.0%$$'

coverage-html:
	go test $(DOMAIN_TEST_PKGS) -coverpkg=$(DOMAIN_COVER_PKG) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
