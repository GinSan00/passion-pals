package notify

import (
	"context"
	"fmt"
	"hunt/internal/repository"
	"log/slog"
	"net/http"

	models "hunt/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type NotifyService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func New(log *slog.Logger, repo *repository.Repository) *NotifyService {
	return &NotifyService{
		log:  log,
		repo: repo,
	}
}

func (notify *NotifyService) GetNotifications(c *gin.Context) {
	// Извлекаем user_id из контекста
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User claims not found"})
		return
	}

	// Приводим к map[string]interface{} сначала
	userClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	// Теперь можно безопасно извлекать данные
	userIDFloat, ok := userClaims["user_id"].(float64)
	userID := int(userIDFloat)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
		return
	}

	// Используем userID для получения профиля
	notifications, err := notify.repo.GetNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func AddNotificaion(ctx context.Context, repo *repository.Repository, userID int, message string, notificationType models.NotificationType) error {

	err := repo.AddNotification(ctx, userID, message, notificationType)
	if err != nil {
		return fmt.Errorf("failed to add notification: %w", err)
	}

	return nil
}
