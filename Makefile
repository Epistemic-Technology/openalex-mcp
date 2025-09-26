# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Binary directory
BINARY_DIR=bin

# Binary names
BINARIES=openalex-mcp-local-server

# Build flags
LDFLAGS=-ldflags "-s -w"

export ORCID_MCP_ROOT := $(CURDIR)

# Default target
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: $(BINARIES)

# Individual binary targets
.PHONY: openalex-mcp-local-server
openalex-mcp-local-server:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/openalex-mcp-local-server ./cmd/openalex-mcp-local-server

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

# Run the server (development)
.PHONY: run
run:
	$(GOCMD) run ./cmd/openalex-mcp-local-server/main.go

# Install binaries to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install ./cmd/...

# Run the MCP inspector on local server
.PHONY: inspect
inspect:
	npx @modelcontextprotocol/inspector $(ORCID_MCP_ROOT)/$(BINARY_DIR)/openalex-mcp-local-server

# Add local server to claude code
.PHONY: cc-add-mcp
cc-add-mcp:
	claude mcp add openalex-mcp-local-server --scope project -- $(ORCID_MCP_ROOT)/$(BINARY_DIR)/openalex-mcp-local-server

# Remove all mcp servers from claude code
.PHONY: cc-remove-mcp
cc-remove-mcp:
	-claude mcp remove openalex-mcp-local-server --scope project

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all                   - Build all binaries (default)"
	@echo "  build                 - Build all binaries"
	@echo "  orcid-mcp-local-server - Build orcid-mcp-local-server binary"
	@echo "  test                  - Run tests"
	@echo "  clean                 - Remove build artifacts"
	@echo "  run                   - Run the server in development mode"
	@echo "  install               - Install binaries to GOPATH/bin"
	@echo "  inspect               - Run the MCP inspector on local server"
	@echo "  cc-add-mcp            - Add local server to claude code"
	@echo "  help                  - Show this help message"
