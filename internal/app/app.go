package app

import (
	"log/slog"
	"time"

	appGrpc "github.com/sy-miller/sso-grpc/internal/app/grpc"
	"github.com/sy-miller/sso-grpc/internal/services/auth"
	"github.com/sy-miller/sso-grpc/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *appGrpc.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTtl time.Duration) *App {
	// Initialise storage
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic("could not initialise storage: " + err.Error())
	}

	// Initialise auth service
	authService := auth.New(log, storage, storage, storage, tokenTtl)

	//Initialise grpcApp

	grpcApp := appGrpc.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}

func (a *App) MustRun() {
	if err := a.GRPCSrv.Run(); err != nil {
		panic(err.Error())
	}
}
