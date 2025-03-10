package profile

import (
	"log/slog"
	"net/http"
	"passion-pals-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ProfileService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func New(log *slog.Logger, repo *repository.Repository) *ProfileService {
	return &ProfileService{
		log:  log,
		repo: repo,
	}
}

func (profile *ProfileService) GetUserProfile(c *gin.Context) {
	// Извлекаем user_id из контекста
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User claims not found"})
		return
	}

	// Приводим к map[string]interface{} сначала
	userClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " Invalid claims format"})
		return
	}

	// Теперь можно безопасно извлекать данные
	userIDFloat, ok := userClaims["user_id"].(float64)
	userID := int(userIDFloat)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " User ID not found in token"})
		return
	}

	// Используем userID для получения профиля
	userProfile, err := profile.repo.GetProfileByUserId(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

func (profile *ProfileService) GetProfileByID(c *gin.Context) { /*
		// Извлекаем user_id из контекста
		userID := c.Param("id")

		userProfile, err := profile.repo.GetProfileByUserId(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
			return
		}

		c.JSON(http.StatusOK, userProfile) */
}

func (profile *ProfileService) GetProfiles(c *gin.Context) {
	// Получаем все анкеты пользователей
	profiles, err := profile.repo.GetProfiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

func (profile *ProfileService) EditUserProfile(c *gin.Context) {
	// Извлекаем user_id из контекста
	/*
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			profile.log.Debug("User ID not found in token")
			return
		}

		// Преобразуем userID в string
		userIdStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			profile.log.Debug("Invalid user ID format")
			return
		}

		// Используем userID для получения профиля
		userProfile, err := profile.repo.GetProfileByUserId(c.Request.Context(), userIdStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
			return
		}

		c.JSON(http.StatusOK, userProfile)
	*/
}

func (profile *ProfileService) DeleteUserProfile(c *gin.Context) {
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

	// Используем userID для удаления профиля
	err := profile.repo.DeleteUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User is sucessifuly deleted"})
}
