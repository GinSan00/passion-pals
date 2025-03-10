package profilehttp

import (
	"hunt/internal/utils/middleware"

	"github.com/gin-gonic/gin"
)

// Profile определяет интерфейс для работы с профилями
type Response interface {
	PostResponse(c *gin.Context)    // Получение профиля текущего пользователя
	GetResponses(c *gin.Context)    // Получение списка всех профилей
	ConfirmResponse(c *gin.Context) // Получение списка всех профилей
	RejectResponse(c *gin.Context)  // Получение списка всех профилей

}

// Register регистрирует маршруты для работы с профилями
func Register(router *gin.Engine, responseService Response) {
	// Группа маршрутов для работы с откликами текущего пользователя
	profileGroup := router.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware()) // Применяем middleware для аутентификации
	{
		// GET /profile/responses - получение откликов текущего пользователя
		profileGroup.GET("/responses", responseService.GetResponses)
		profileGroup.PUT("/responses/:id", responseService.ConfirmResponse)
		profileGroup.DELETE("/responses/:id", responseService.RejectResponse)
	}

	profilesGroup := router.Group("/profiles/:id")
	profilesGroup.Use(middleware.AuthMiddleware()) // Применяем middleware для аутентификации
	{
		// POST /profiles/:id - добавление нового отклика на указанный :id
		profilesGroup.POST("", responseService.PostResponse)
	}
}
