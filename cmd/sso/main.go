package main

import (
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"sso-grpc-ntc/internal/config"
	"sso-grpc-ntc/pkg/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Запуск приложения")

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сервер приложения
}

// setupLogger Создаёт логгера в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// setupPrettySlog Создает логгера для локального окружения
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
