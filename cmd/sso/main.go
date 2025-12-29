package main

import (
	"github.com/sy-miller/sso-grpc/internal/config"
)

func main() {
	// Initialise configuration
	cfg := config.MustLoad()

	// Initialise logger
	_ = cfg.LogLevel

	// Initialise 'app'

	// Start gRPC server
}
