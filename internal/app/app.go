package app

import (
	"log/slog"
	grpcapp "sso-grpc-ntc/internal/app/grpc"
	"sso-grpc-ntc/internal/services/auth"
	"sso-grpc-ntc/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCServer: grpcApp}
}
