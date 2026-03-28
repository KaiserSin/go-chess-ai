SCRIPT_DIR := ./scripts

.PHONY: help test coverage coverage-html lint run

help:
	@printf '%s\n' \
		'Available targets:' \
		'  make test           Run all Go tests' \
		'  make coverage       Run domain coverage checks and enforce >= 90.0%' \
		'  make coverage-html  Generate coverage.out and coverage.html' \
		'  make lint           Run go vet and golangci-lint' \
		'  make run            Launch the desktop chess app'

test:
	@$(SCRIPT_DIR)/test.sh

coverage:
	@$(SCRIPT_DIR)/coverage.sh

coverage-html:
	@$(SCRIPT_DIR)/coverage.sh --html

lint:
	@$(SCRIPT_DIR)/lint.sh

run:
	@$(SCRIPT_DIR)/run-chess.sh
