package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/mrauhala/mcp-ogc-edr/internal/mcp"
)

func main() {
	var (
		edrBaseURL = flag.String("edr-url", envOrDefault("EDR_BASE_URL", "https://example.com/edr"), "OGC EDR API base URL")
		transport  = flag.String("transport", envOrDefault("MCP_TRANSPORT", "stdio"), "MCP transport: stdio, sse, or streamable-http")
		sseAddr    = flag.String("sse-addr", envOrDefault("SSE_ADDR", ":8080"), "SSE listen address")
		logLevel   = flag.String("log-level", envOrDefault("LOG_LEVEL", "info"), "Log level: debug, info, warn, error")
	)
	flag.Parse()

	// Setup structured logger
	var level slog.Level
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	slog.Info("starting MCP OGC EDR server",
		"edr_url", *edrBaseURL,
		"transport", *transport,
	)

	srv, err := mcp.NewServer(*edrBaseURL)
	if err != nil {
		slog.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := srv.Run(ctx, *transport, *sseAddr); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
