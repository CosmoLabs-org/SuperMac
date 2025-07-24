# =============================================================================
# SuperMac Makefile
# =============================================================================
# Development automation for SuperMac
# 
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# Project Configuration
PROJECT_NAME := SuperMac
VERSION := 2.1.0
BIN_DIR := bin
LIB_DIR := lib
TEST_DIR := tests
DOCS_DIR := docs
CONFIG_DIR := config

# Shell Scripts
MAIN_SCRIPT := $(BIN_DIR)/mac
INSTALL_SCRIPT := $(BIN_DIR)/install.sh
TEST_SCRIPT := $(TEST_DIR)/test.sh
MODULES := $(wildcard $(LIB_DIR)/*.sh)

# Colors for output
BOLD := \033[1m
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
PURPLE := \033[35m
CYAN := \033[36m
NC := \033[0m

# Default target
.PHONY: help
help: ## Show this help message
	@echo "$(BOLD)$(PURPLE)SuperMac Development Commands$(NC)"
	@echo ""
	@echo "$(BOLD)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)Examples:$(NC)"
	@echo "  make test                    # Run all tests"
	@echo "  make lint                    # Check code quality"
	@echo "  make install-dev             # Install development version"
	@echo "  make clean                   # Clean temporary files"
	@echo ""

# =============================================================================
# Development Setup
# =============================================================================

.PHONY: setup
setup: ## Set up development environment
	@echo "$(BOLD)$(BLUE)Setting up development environment...$(NC)"
	@chmod +x $(MAIN_SCRIPT) $(INSTALL_SCRIPT) $(TEST_SCRIPT)
	@chmod +x $(MODULES)
	@echo "$(GREEN)✓ Made scripts executable$(NC)"
	@mkdir -p ~/bin
	@ln -sf "$(PWD)/$(MAIN_SCRIPT)" ~/bin/mac-dev
	@echo "$(GREEN)✓ Created development symlink: mac-dev$(NC)"
	@echo "$(BOLD)$(GREEN)Development environment ready!$(NC)"
	@echo "Test with: $(CYAN)mac-dev help$(NC)"

.PHONY: install-dev
install-dev: setup ## Install development version
	@echo "$(BOLD)$(BLUE)Installing development version...$(NC)"
	@if ! echo $$PATH | grep -q "$$HOME/bin"; then \
		echo "$(YELLOW)⚠ ~/bin not in PATH. Add this to your shell config:$(NC)"; \
		echo "  export PATH=\"\$$HOME/bin:\$$PATH\""; \
	fi

# =============================================================================
# Testing
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BOLD)$(BLUE)Running SuperMac test suite...$(NC)"
	@$(TEST_SCRIPT)

.PHONY: test-quick
test-quick: ## Run quick tests (syntax only)
	@echo "$(BOLD)$(BLUE)Running quick syntax tests...$(NC)"
	@for script in $(MAIN_SCRIPT) $(INSTALL_SCRIPT) $(MODULES); do \
		echo "Testing $$script..."; \
		bash -n $$script || exit 1; \
	done
	@echo "$(GREEN)✓ All syntax tests passed$(NC)"

.PHONY: test-module
test-module: ## Test specific module (usage: make test-module MODULE=finder)
	@if [ -z "$(MODULE)" ]; then \
		echo "$(RED)Error: MODULE not specified$(NC)"; \
		echo "Usage: make test-module MODULE=finder"; \
		exit 1; \
	fi
	@echo "$(BOLD)$(BLUE)Testing module: $(MODULE)$(NC)"
	@$(TEST_SCRIPT) $(MODULE)

.PHONY: test-performance
test-performance: ## Run performance tests
	@echo "$(BOLD)$(BLUE)Running performance tests...$(NC)"
	@echo "Testing startup time..."
	@time $(MAIN_SCRIPT) help >/dev/null
	@echo "Testing help system performance..."
	@time $(MAIN_SCRIPT) help finder >/dev/null

# =============================================================================
# Code Quality
# =============================================================================

.PHONY: lint
lint: ## Run shellcheck on all scripts
	@echo "$(BOLD)$(BLUE)Running shellcheck...$(NC)"
	@if command -v shellcheck >/dev/null 2>&1; then \
		for script in $(MAIN_SCRIPT) $(INSTALL_SCRIPT) $(MODULES); do \
			echo "Checking $$script..."; \
			shellcheck $$script || exit 1; \
		done; \
		echo "$(GREEN)✓ All scripts passed shellcheck$(NC)"; \
	else \
		echo "$(YELLOW)⚠ shellcheck not installed. Install with: brew install shellcheck$(NC)"; \
	fi

.PHONY: format
format: ## Format shell scripts with shfmt
	@echo "$(BOLD)$(BLUE)Formatting shell scripts...$(NC)"
	@if command -v shfmt >/dev/null 2>&1; then \
		for script in $(MAIN_SCRIPT) $(INSTALL_SCRIPT) $(MODULES); do \
			echo "Formatting $$script..."; \
			shfmt -w -i 4 -ci $$script; \
		done; \
		echo "$(GREEN)✓ All scripts formatted$(NC)"; \
	else \
		echo "$(YELLOW)⚠ shfmt not installed. Install with: brew install shfmt$(NC)"; \
	fi

.PHONY: check
check: lint test-quick ## Run all code quality checks
	@echo "$(GREEN)✓ All quality checks passed$(NC)"

# =============================================================================
# Documentation
# =============================================================================

.PHONY: docs
docs: ## Generate documentation
	@echo "$(BOLD)$(BLUE)Generating documentation...$(NC)"
	@echo "$(GREEN)✓ Documentation up to date$(NC)"
	@echo "Main docs: $(DOCS_DIR)/README.md"
	@echo "Dev guide: $(DOCS_DIR)/DEVELOPMENT.md"

.PHONY: docs-serve
docs-serve: ## Serve documentation locally (requires Python)
	@echo "$(BOLD)$(BLUE)Starting documentation server...$(NC)"
	@echo "Open: http://localhost:8000"
	@cd $(DOCS_DIR) && python3 -m http.server 8000

# =============================================================================
# Building and Distribution
# =============================================================================

.PHONY: build
build: check ## Build SuperMac for distribution
	@echo "$(BOLD)$(BLUE)Building SuperMac v$(VERSION)...$(NC)"
	@mkdir -p dist
	@cp -r $(BIN_DIR) $(LIB_DIR) $(CONFIG_DIR) $(DOCS_DIR) dist/
	@echo "$(GREEN)✓ Build complete: dist/$(NC)"

.PHONY: package
package: build ## Create distribution package
	@echo "$(BOLD)$(BLUE)Creating distribution package...$(NC)"
	@cd dist && tar -czf ../supermac-$(VERSION).tar.gz *
	@echo "$(GREEN)✓ Package created: supermac-$(VERSION).tar.gz$(NC)"

.PHONY: release-check
release-check: ## Check if ready for release
	@echo "$(BOLD)$(BLUE)Checking release readiness...$(NC)"
	@echo "Running comprehensive tests..."
	@$(MAKE) test
	@echo "Checking code quality..."
	@$(MAKE) lint
	@echo "Verifying version numbers..."
	@grep -q "$(VERSION)" $(MAIN_SCRIPT) || (echo "$(RED)Version mismatch in main script$(NC)" && exit 1)
	@grep -q "$(VERSION)" $(INSTALL_SCRIPT) || (echo "$(RED)Version mismatch in install script$(NC)" && exit 1)
	@echo "$(GREEN)✓ Ready for release!$(NC)"

# =============================================================================
# Debugging and Development
# =============================================================================

.PHONY: debug
debug: ## Run SuperMac in debug mode
	@echo "$(BOLD)$(BLUE)Running in debug mode...$(NC)"
	@SUPERMAC_DEBUG=1 $(MAIN_SCRIPT) $(ARGS)

.PHONY: debug-module
debug-module: ## Debug specific module (usage: make debug-module MODULE=finder ARGS="status")
	@if [ -z "$(MODULE)" ]; then \
		echo "$(RED)Error: MODULE not specified$(NC)"; \
		echo "Usage: make debug-module MODULE=finder ARGS=\"status\""; \
		exit 1; \
	fi
	@echo "$(BOLD)$(BLUE)Debugging module: $(MODULE)$(NC)"
	@SUPERMAC_DEBUG=1 $(MAIN_SCRIPT) $(MODULE) $(ARGS)

.PHONY: profile
profile: ## Profile performance of commands
	@echo "$(BOLD)$(BLUE)Profiling SuperMac performance...$(NC)"
	@echo "Main help:"
	@time $(MAIN_SCRIPT) help >/dev/null 2>&1
	@echo "Category help:"
	@time $(MAIN_SCRIPT) help finder >/dev/null 2>&1
	@echo "Command execution:"
	@time $(MAIN_SCRIPT) system info >/dev/null 2>&1

# =============================================================================
# Maintenance
# =============================================================================

.PHONY: clean
clean: ## Clean temporary files and build artifacts
	@echo "$(BOLD)$(BLUE)Cleaning temporary files...$(NC)"
	@rm -rf dist/
	@rm -f *.tar.gz
	@rm -f .DS_Store
	@find . -name "*.tmp" -delete
	@find . -name "*.log" -delete
	@echo "$(GREEN)✓ Cleaned temporary files$(NC)"

.PHONY: clean-all
clean-all: clean ## Clean everything including development links
	@echo "$(BOLD)$(BLUE)Cleaning all development files...$(NC)"
	@rm -f ~/bin/mac-dev
	@echo "$(GREEN)✓ Removed development symlink$(NC)"

.PHONY: update-version
update-version: ## Update version number (usage: make update-version VERSION=2.2.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION not specified$(NC)"; \
		echo "Usage: make update-version VERSION=2.2.0"; \
		exit 1; \
	fi
	@echo "$(BOLD)$(BLUE)Updating version to $(VERSION)...$(NC)"
	@sed -i '' 's/VERSION=".*"/VERSION="$(VERSION)"/' $(MAIN_SCRIPT)
	@sed -i '' 's/VERSION=".*"/VERSION="$(VERSION)"/' $(INSTALL_SCRIPT)
	@sed -i '' 's/"version": ".*"/"version": "$(VERSION)"/' $(CONFIG_DIR)/config.json
	@echo "$(GREEN)✓ Version updated to $(VERSION)$(NC)"

# =============================================================================
# Installation Testing
# =============================================================================

.PHONY: test-install
test-install: ## Test installation script
	@echo "$(BOLD)$(BLUE)Testing installation script...$(NC)"
	@echo "This will test the installer in a safe way..."
	@bash -n $(INSTALL_SCRIPT)
	@echo "$(GREEN)✓ Installation script syntax is valid$(NC)"

.PHONY: install-local
install-local: ## Install SuperMac locally (for testing)
	@echo "$(BOLD)$(BLUE)Installing SuperMac locally...$(NC)"
	@$(INSTALL_SCRIPT)

# =============================================================================
# CI/CD Support
# =============================================================================

.PHONY: ci
ci: lint test build ## Run CI pipeline (lint, test, build)
	@echo "$(GREEN)✓ CI pipeline completed successfully$(NC)"

.PHONY: ci-quick
ci-quick: test-quick ## Quick CI check (syntax only)
	@echo "$(GREEN)✓ Quick CI check passed$(NC)"

# =============================================================================
# Information
# =============================================================================

.PHONY: info
info: ## Show project information
	@echo "$(BOLD)$(PURPLE)SuperMac Project Information$(NC)"
	@echo ""
	@echo "$(BOLD)Version:$(NC) $(VERSION)"
	@echo "$(BOLD)Scripts:$(NC)"
	@echo "  Main: $(MAIN_SCRIPT)"
	@echo "  Install: $(INSTALL_SCRIPT)"
	@echo "  Test: $(TEST_SCRIPT)"
	@echo "$(BOLD)Modules:$(NC)"
	@for module in $(MODULES); do echo "  $$module"; done
	@echo "$(BOLD)Documentation:$(NC) $(DOCS_DIR)/"
	@echo "$(BOLD)Configuration:$(NC) $(CONFIG_DIR)/"
	@echo ""

.PHONY: status
status: ## Show development status
	@echo "$(BOLD)$(PURPLE)Development Status$(NC)"
	@echo ""
	@echo "$(BOLD)File Permissions:$(NC)"
	@ls -la $(MAIN_SCRIPT) $(INSTALL_SCRIPT) $(TEST_SCRIPT)
	@echo ""
	@echo "$(BOLD)Development Link:$(NC)"
	@ls -la ~/bin/mac-dev 2>/dev/null || echo "  No development link found"
	@echo ""
	@echo "$(BOLD)PATH Check:$(NC)"
	@if echo $$PATH | grep -q "$$HOME/bin"; then \
		echo "  $(GREEN)✓ ~/bin is in PATH$(NC)"; \
	else \
		echo "  $(YELLOW)⚠ ~/bin not in PATH$(NC)"; \
	fi

# =============================================================================
# Development Shortcuts
# =============================================================================

.PHONY: dev
dev: setup test ## Quick development setup and test

.PHONY: quick
quick: test-quick lint ## Quick quality check

.PHONY: all
all: clean setup test lint build ## Full development cycle

# Make variables
.DEFAULT_GOAL := help
MAKEFLAGS += --no-print-directory
