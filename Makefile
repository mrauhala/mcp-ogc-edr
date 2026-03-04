BINARY := mcp-ogc-edr
MODULE  := github.com/mrauhala/mcp-ogc-edr

.PHONY: all build run tidy lint test docker-build clean

all: build

build:
	go build -o bin/$(BINARY) ./cmd/server

run: build
	EDR_BASE_URL=https://your-edr-server.com/edr ./bin/$(BINARY)

run-sse: build
	EDR_BASE_URL=https://your-edr-server.com/edr \
	MCP_TRANSPORT=sse \
	SSE_ADDR=:8080 \
	./bin/$(BINARY)

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

test:
	go test ./...

docker-build:
	docker build -t $(BINARY):latest .

clean:
	rm -rf bin/
