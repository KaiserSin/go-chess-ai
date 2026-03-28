SCRIPT_DIR := ./scripts

.PHONY: help test coverage coverage-html run

help:
	@printf '%s\n' \
		'Available targets:' \
		'  make test           Run all Go tests' \
		'  make coverage       Run domain coverage checks and enforce 100.0%' \
		'  make coverage-html  Generate coverage.out and coverage.html' \
		'  make run            Launch the desktop chess app'

test:
	@$(SCRIPT_DIR)/test.sh

coverage:
	@$(SCRIPT_DIR)/coverage.sh

coverage-html:
	@$(SCRIPT_DIR)/coverage.sh --html

run:
	@$(SCRIPT_DIR)/run-chess.sh
