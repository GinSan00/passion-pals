package profilehttp

import (
	"passion-pals-backend/internal/utils/middleware"

	"github.com/gin-gonic/gin"
)

// Profile определяет интерфейс для работы с профилями
type Profile interface {
	GetUserProfile(c *gin.Context)    // Получение профиля текущего пользователя
	GetProfiles(c *gin.Context)       // Получение списка всех профилей
	GetProfileByID(c *gin.Context)    // Получение профиля по ID
	EditUserProfile(c *gin.Context)   // Редактирование профиля текущего пользователя
	DeleteUserProfile(c *gin.Context) // Редактирование профиля текущего пользователя
}

// Register регистрирует маршруты для работы с профилями
func Register(router *gin.Engine, profileService Profile) {
	// Группа маршрутов для работы с профилем текущего пользователя
	profileGroup := router.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware()) // Применяем middleware для аутентификации
	{
		// GET /profile - получение профиля текущего пользователя
		profileGroup.GET("", profileService.GetUserProfile)

		// PUT /profile - редактирование профиля текущего пользователя
		profileGroup.PUT("", profileService.EditUserProfile)

		profileGroup.DELETE("", profileService.DeleteUserProfile)
	}

	// Группа маршрутов для работы с профилями других пользователей
	profilesGroup := router.Group("/profiles")
	profilesGroup.Use(middleware.AuthMiddleware()) // Применяем middleware для аутентификации
	{
		// GET /profiles - получение списка всех профилей
		profilesGroup.GET("", profileService.GetProfiles)

		// GET /profiles/:id - получение профиля по ID
		profilesGroup.GET("/:id", profileService.GetProfileByID)
	}
}
