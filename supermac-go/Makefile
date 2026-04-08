BINARY = mac
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test lint install clean

build:
	go build -ldflags "-X github.com/cosmolabs-org/supermac/internal/version.Version=$(VERSION) -X github.com/cosmolabs-org/supermac/internal/version.BuildDate=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(BINARY) ./cmd/mac

test:
	go test ./...

lint:
	golangci-lint run

install: build
	cp $(BINARY) ~/bin/

clean:
	rm -f $(BINARY)
