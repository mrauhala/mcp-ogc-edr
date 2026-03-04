FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/mcp-ogc-edr ./cmd/server

FROM scratch
COPY --from=builder /bin/mcp-ogc-edr /mcp-ogc-edr
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV EDR_BASE_URL=""
ENV MCP_TRANSPORT=stdio
ENV LOG_LEVEL=info

ENTRYPOINT ["/mcp-ogc-edr"]
