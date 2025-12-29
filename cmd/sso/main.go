package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/sy-miller/sso-grpc/internal/config"
)

var log *slog.Logger

func main() {
	// Initialise configuration
	cfg := config.MustLoad()

	// Initialise logger
	log = setupLogger(cfg.LogLevel).With(
		slog.String("environment", cfg.Env),
	)

	log.Info("starting application",
		slog.Int("port", *cfg.GRPC.Port),
		slog.Any("timeout", cfg.GRPC.Timeout),
	)

	// Initialise 'app'

	// Start gRPC server
}

func setupLogger(logLevel string) *slog.Logger {
	var (
		log       *slog.Logger
		slogLevel slog.Level
	)

	switch logLevel {
	case config.DEBUG:
		slogLevel = slog.LevelDebug
	case config.WARN:
		slogLevel = slog.LevelWarn
	case config.ERROR:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Log UTC time
				if a.Key == "time" {
					a.Value = slog.AnyValue(time.Now().UTC())
				}
				return a
			},
		}),
	)

	return log
}
