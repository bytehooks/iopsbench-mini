.PHONY: build test clean version

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X iopsbench-mini/internal/version.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o iopsbench-mini

test:
	go test ./...

clean:
	rm -f iopsbench-mini iopsbench-mini-*

version:
	@echo $(VERSION)
