BINARY  := rmx
BIN_DIR := ./bin
MODULE  := github.com/braswelljr/rmx
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X $(MODULE)/internal/common.Version=$(VERSION)"

HASGOCILINT := $(shell which golangci-lint 2> /dev/null)
ifdef HASGOCILINT
    GOLINT = golangci-lint
else
    GOLINT = bin/golangci-lint
endif

.PHONY: help
help: ## List available targets
	@grep -hE '^[a-zA-Z0-9_/.-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ── Build ───────────────────────────────────────────────────────────────────
.PHONY: build
build: ## Build the binary for the host platform into bin/
	go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY) .

.PHONY: build/all
build/all: ## Cross-compile for linux/darwin/windows (amd64+arm64)
	GOOS=linux   GOARCH=amd64 go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY)-windows-amd64.exe .

.PHONY: install
install: ## Install rmx with version ldflags
	go install -trimpath $(LDFLAGS) .

.PHONY: run
run: ## Run rmx from source (pass args with ARGS=...)
	go run . $(ARGS)

# ── Dependencies ──────────────────────────────────────────────────────────────
.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: download
download: ## Download module dependencies
	go mod download

# ── Quality ─────────────────────────────────────────────────────────────────
.PHONY: test
test: ## Run tests with the race detector
	go test -race ./...

.PHONY: cover
cover: ## Run tests and write coverage.out
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	$(GOLINT) run

.PHONY: fix
fix: ## Format and tidy imports
	gofmt -s -w .
	goimports -w $$(find . -type f -name '*.go' -not -path "*/vendor/*")

.PHONY: ci
ci: vet lint test docs/check ## Run the full CI suite locally

# ── Documentation ─────────────────────────────────────────────────────────────
.PHONY: docs
docs: ## Regenerate the CLI reference (all formats) into docs/
	go run ./tools/gendocs docs all

.PHONY: docs/check
docs/check: docs ## Fail if the committed docs are out of date
	@if [ -n "$$(git status --porcelain docs)" ]; then \
		echo "docs/ is out of date — commit the regenerated files"; \
		git status --porcelain docs; \
		exit 1; \
	fi

# ── Docker ────────────────────────────────────────────────────────────────────
.PHONY: docker/build
docker/build: ## Build the Docker image (tag: rmx:$(VERSION))
	docker build --build-arg VERSION=$(VERSION) -t $(BINARY):$(VERSION) -t $(BINARY):latest .

.PHONY: docker/run
docker/run: ## Run rmx in Docker against the current directory (ARGS=...)
	docker run --rm -v "$$(pwd):/work" $(BINARY):latest $(ARGS)

# ── Clean ─────────────────────────────────────────────────────────────────────
.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) dist coverage.out
