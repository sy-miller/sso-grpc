package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sy-miller/sso-grpc/internal/app"
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
	application := app.New(log, *cfg.GRPC.Port, *cfg.StoragePath, time.Duration(cfg.TokenTTL))

	// Start gRPC server to listen in a goroutine
	go application.MustRun()

	// Listen for signals in the main goroutine
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := <-stop

	log.Info("stopping application", slog.String("signal", s.String()))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
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
