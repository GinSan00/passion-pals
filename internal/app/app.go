package app

import (
	httppapp "hunt/internal/app/httpapp"
	"hunt/internal/controllers/auth"
	"hunt/internal/controllers/profile"
	"hunt/internal/repository"
	"log/slog"
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
