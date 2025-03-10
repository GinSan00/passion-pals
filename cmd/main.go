package main

import (
	"log/slog"
	"os"
	"os/signal"
	"passion-pals-backend/internal/config"
	"syscall"

	"passion-pals-backend/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadByPath("./config/local.yaml")

	log := setupLogger(cfg.Env)

	log.Info("Starting application", slog.Any("cfg", cfg))

	application := app.New(log, cfg.Server.Port, cfg.ConnectionString, cfg.TokenTTL)

	go application.HTTPSrv.MustRun()

	//Мягкое завершение
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.HTTPSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
