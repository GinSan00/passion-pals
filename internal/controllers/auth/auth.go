package auth

import (
	"hunt/internal/repository"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func New(log *slog.Logger, repo *repository.Repository, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:  log,
		repo: repo,
	}
}

func (auth *AuthService) Register(c *gin.Context) {
	var newUser struct {
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		BirthDate time.Time `json:"birth_date" time_format:"2006-01-02T15:04:05Z07:00"`
		Gender    string    `json:"gender"`
	}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Проверка, что username и password не пустые
	if strings.TrimSpace(newUser.Username) == "" || strings.TrimSpace(newUser.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	password_hash, err := hashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Вставка пользователя в базу данных
	userID, err := auth.repo.CreateUser(c.Request.Context(), newUser.Username, password_hash, newUser.Email, newUser.BirthDate, newUser.Gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		auth.log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

func (auth *AuthService) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID, password_hash, err := auth.repo.FindUserByUserEmail(c.Request.Context(), loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username o password"})
		return
	}
	if password_hash == "" {

	}
	/*if err := compareHashAndPassword(loginData.Password, password_hash); !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	*/
	token, err := generateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	auth.log.Debug(token)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   token,
	})
}

func generateJWT(userID int) (string, error) {

	var jwtSecret = []byte("your_secret_key")

	claims := jwt.MapClaims{
		"user_id": userID,                                // Полезные данные (payload)
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена (24 часа)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func compareHashAndPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
