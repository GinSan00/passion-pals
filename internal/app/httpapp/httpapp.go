package httppapp

import (
	"fmt"
	authhttp "hunt/internal/http/auth" // Предположим, что у вас есть HTTP-хендлеры для auth
	profilehttp "hunt/internal/http/profile"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type App struct {
	log    *slog.Logger
	router *gin.Engine
	port   int
}

func New(
	log *slog.Logger,
	authService authhttp.Auth,
	profileService profilehttp.Profile, // Предположим, что у вас есть HTTP-хендлер для auth
	port int,
) *App {
	// Инициализация Gin
	router := gin.Default()

	//Настройка разрешенных источников и методов запроса TODO - вынести в отдельный метод
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	// Регистрация HTTP-хендлеров
	authhttp.Register(router, authService)
	profilehttp.Register(router, profileService)

	return &App{
		log:    log,
		router: router,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httppapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port))

	// Создание HTTP-сервера
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: a.router,
	}

	log.Info("HTTP server is running", slog.String("addr", server.Addr))

	// Запуск HTTP-сервера
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "httppapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping HTTP server", slog.Int("port", a.port))

	// Здесь можно добавить логику для graceful shutdown, если необходимо
	// Например, использование context.WithTimeout и server.Shutdown
}
