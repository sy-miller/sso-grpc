package app

import (
	"log/slog"
	"time"

	appGrpc "github.com/sy-miller/sso-grpc/internal/app/grpc"
)

type App struct {
	GRPCSrv *appGrpc.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTtl time.Duration) *App {
	// Initialise storag

	// Initialise auth service

	//Initialise grpcApp

	grpcApp := appGrpc.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}

func (a *App) MustRun() {
	if err := a.GRPCSrv.Run(); err != nil {
		panic(err.Error())
	}
}
