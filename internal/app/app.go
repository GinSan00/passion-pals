package app

import (
	"log/slog"
	httppapp "passion-pals-backend/internal/app/httpapp"
	"passion-pals-backend/internal/controllers/auth"
	"passion-pals-backend/internal/controllers/profile"
	"passion-pals-backend/internal/repository"
	"time"
)

type App struct {
	HTTPSrv *httppapp.App
}

func New(
	log *slog.Logger,
	httpPort int,
	connStr string,
	tokenTTL time.Duration,
) *App {

	repo, err := repository.NewRepository(connStr)

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, repo, tokenTTL)
	profileService := profile.New(log, repo)

	httpApp := httppapp.New(log, authService, profileService, httpPort)

	return &App{
		HTTPSrv: httpApp,
	}
}
