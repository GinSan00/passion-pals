package profilehttp

import (
	"passion-pals-backend/internal/utils/middleware"

	"github.com/gin-gonic/gin"
)

// Profile определяет интерфейс для работы с профилями
type Notification interface {
	GetNotifications(c *gin.Context) // Получение профиля текущего пользователя
	MarkAsRead(c *gin.Context)       // Получение списка всех профилей
}

// Register регистрирует маршруты для работы с профилями
func Register(router *gin.Engine, notificationService Notification) {
	// Группа маршрутов для работы с уведомлениями текущего пользователя
	profileGroup := router.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware())
	{
		profileGroup.GET("/notifications", notificationService.GetNotifications)    // Получить уведомления
		profileGroup.PUT("/notifications/:id/read", notificationService.MarkAsRead) // Отметить как прочитанное
	}
}
