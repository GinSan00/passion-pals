package responses

import (
	"hunt/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type ResponsesService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func New(log *slog.Logger, repo *repository.Repository) *ResponsesService {
	return &ResponsesService{
		log:  log,
		repo: repo,
	}
}

func (response *ResponsesService) PostResponse(c *gin.Context) {
	userID := c.Param("id")
	profileID := c.Param("profileID")

	err := response.repo.AddResponse(c.Request.Context(), userID, profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}
}

func (response *ResponsesService) ConfirmResponse(c *gin.Context) {
	responseID := c.Param("id")

	err := response.repo.ConfirmResponse(c.Request.Context(), responseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}
}

func (response *ResponsesService) RejectResponse(c *gin.Context) {
	responseID := c.Param("id")

	err := response.repo.RejectResponse(c.Request.Context(), responseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}
}

func (response *ResponsesService) GetResponses(c *gin.Context) {
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
	incomingResponses, err := response.repo.GetIncomingResponses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	outgoingResponses, err := response.repo.GetOutgoingResponses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	// Создаем структуру для объединения ответов
	responseData := gin.H{
		"incoming_responses": incomingResponses,
		"outgoing_responses": outgoingResponses,
	}

	// Отправляем объединенные данные
	c.JSON(http.StatusOK, responseData)
}
