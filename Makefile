SCRIPT_DIR := ./scripts

.PHONY: help test test-ai-extended coverage coverage-html lint run

help:
	@printf '%s\n' \
		'Available targets:' \
		'  make test           Run all Go tests' \
		'  make test-ai-extended  Run the extended deterministic AI test suite' \
		'  make coverage       Run core-package coverage checks and enforce >= 75.0%' \
		'  make coverage-html  Generate coverage.out and coverage.html' \
		'  make lint           Run go vet and golangci-lint' \
		'  make run            Launch the desktop chess app'

test:
	@$(SCRIPT_DIR)/test.sh

test-ai-extended:
	@sh $(SCRIPT_DIR)/test-ai-extended.sh

coverage:
	@$(SCRIPT_DIR)/coverage.sh

coverage-html:
	@$(SCRIPT_DIR)/coverage.sh --html

lint:
	@$(SCRIPT_DIR)/lint.sh

run:
	@$(SCRIPT_DIR)/run-chess.sh
