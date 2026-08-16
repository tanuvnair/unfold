# unfold — local CLI, API, and web UI
#
# Usage:  make <target>
# Config: copy .env.example → .env (optional); make includes it when present.

.DEFAULT_GOAL := help

DIST        ?= dist
BIN         ?= unfold
API_BIN     ?= unfold-api
GO          ?= go
WEB_DIR     ?= web
NPM         ?= npm
WEB_DIST    := $(WEB_DIR)/dist
WEBUI_DIST  := internal/webui/dist

# Defaults; .env values (when present) override via include/export below.
UNFOLD_API_CONFIG ?= configs/banks.json

# Load repo-root .env into the environment when present.
ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: help \
	cli api api-build \
	web web-install web-build sync-ui \
	test fmt vet tidy check \
	build serve clean

help: ## Show available targets
	@printf 'unfold make targets\n\n'
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@printf '\nDefaults: DIST=$(DIST)  BIN=$(BIN)  API_BIN=$(API_BIN)  UNFOLD_API_CONFIG=$(UNFOLD_API_CONFIG)\n'
	@printf 'Env:      UNFOLD_API_HOST UNFOLD_API_PORT UNFOLD_API_CONFIG UNFOLD_API_CORS_ORIGIN\n'
	@printf '          UNFOLD_WEB_HOST UNFOLD_WEB_PORT UNFOLD_WEB_API_BASE UNFOLD_API_PROXY_TARGET\n'

##@ Go

cli: ## Build the CLI binary → ./dist/unfold
	mkdir -p $(DIST)
	$(GO) build -o $(DIST)/$(BIN) ./cmd/cli

api: ## Run the HTTP API (loads UNFOLD_* from .env)
	$(GO) run ./cmd/api -config $(UNFOLD_API_CONFIG)

sync-ui: web-build ## Copy web/dist into the embed tree for the API binary
	rm -rf $(WEBUI_DIST)
	mkdir -p $(WEBUI_DIST)
	cp -a $(WEB_DIST)/. $(WEBUI_DIST)/

api-build: sync-ui ## Build API+UI binary → ./dist/unfold-api
	mkdir -p $(DIST)
	$(GO) build -o $(DIST)/$(API_BIN) ./cmd/api

serve: api-build ## Build and run the single-binary UI+API
	./$(DIST)/$(API_BIN) -config $(UNFOLD_API_CONFIG)

##@ Web

web-install: ## Install web npm dependencies
	$(NPM) --prefix $(WEB_DIR) install

web: ## Run the Vite web UI (proxies /api to the Go API)
	$(NPM) --prefix $(WEB_DIR) run dev

web-build: ## Build the web UI for production → ./web/dist
	$(NPM) --prefix $(WEB_DIR) run build

##@ Quality

test: ## Run Go tests
	$(GO) test ./...

fmt: ## Format Go sources
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

tidy: ## Tidy go.mod / go.sum
	$(GO) mod tidy

check: fmt vet test ## Format, vet, and test

##@ Meta

build: cli api-build ## Build CLI and API+UI binaries

clean: ## Remove built binaries and restore placeholder embed dist
	rm -rf $(DIST) $(WEB_DIST) $(WEBUI_DIST)
	rm -f $(BIN) $(API_BIN)
	mkdir -p $(WEBUI_DIST)
	touch $(WEBUI_DIST)/.gitkeep
