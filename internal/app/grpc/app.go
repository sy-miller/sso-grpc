package grpc

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/sy-miller/sso-grpc/internal/grpc/auth"
	"google.golang.org/grpc"
)

const pkgFn = "grpcApp"

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const fn = pkgFn + ".Run"
	log := a.log.With(slog.String("fn", fn))

	log.Info("starting gRPC server", slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (a *App) Stop() {
	const fn = pkgFn + ".Stop"

	a.log.With(slog.String("fn", fn)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
