package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	models "hunt/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

// NewRepository создает новый экземпляр Repository
func NewRepository(connStr string) (*Repository, error) {

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// CreateUser создает нового пользователя в базе данных
func (r *Repository) CreateUser(ctx context.Context, username, password, email string, birth_date time.Time, gender string) (int, error) {
	var userID int

	err := r.db.QueryRow(ctx,
		"INSERT INTO users (username, email, passoword, date_of_birth, gender) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		username, email, password, birth_date, gender).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	pErr := r.db.QueryRow(ctx,
		"INSERT INTO profiles (user_id, gender, age, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, gender, сalculateAge(birth_date), time.Now(), time.Now())

	if pErr != nil {
		return userID, fmt.Errorf("failed to create profile: %w", err)
	}

	return userID, nil
}

func сalculateAge(birthDate time.Time) int {
	now := time.Now()

	// Вычисляем разницу в годах
	age := now.Year() - birthDate.Year()

	// Проверяем, был ли уже день рождения в этом году
	if now.YearDay() < birthDate.YearDay() {
		age-- // Если день рождения ещё не наступил, уменьшаем возраст на 1
	}

	return age
}

// FindUserByUsername ищет пользователя по имени пользователя
func (r *Repository) FindUserByUserEmail(ctx context.Context, email string) (int, string, error) {
	var id int
	var password string

	err := r.db.QueryRow(ctx,
		"SELECT id, password FROM users WHERE email = $1",
		email).Scan(&id, &password)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("user not found")
		}
		return 0, "", fmt.Errorf("failed to find user: %w", err)
	}

	return id, password, nil
}

// GetProfileByUserId возвращает данные профиля пользователя по id
func (r *Repository) GetProfileByUserId(ctx context.Context, userId int) (*models.UserProfile, error) {
	var username string
	var age int
	var avatarUrl string
	var aboutMe string
	var gender string
	var lookingFor string
	var createdAt time.Time
	var updatedAt time.Time

	err := r.db.QueryRow(ctx,
		`SELECT 
            u.username, 
            p.age, 
            p.avatar_url, 
            p.about_me, 
			p.gender,
            p.looking_for, 
            p.created_at, 
            p.updated_at 
        FROM 
            profiles p
        JOIN 
            users u ON p.user_id = u.id
        WHERE 
            p.user_id = $1`,
		userId).Scan(&username, &age, &avatarUrl, &aboutMe, &gender, &lookingFor, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &models.UserProfile{
		Username:   username,
		Age:        age,
		AvatarUrl:  avatarUrl,
		AboutMe:    aboutMe,
		Gender:     gender,
		LookingFor: lookingFor,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

// GetProfiles список актуальных анкет
func (r *Repository) GetProfiles(ctx context.Context) ([]*models.UserProfile, error) {
	// Определяем временную границу для последних 7 дней
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)

	// Выполняем запрос к базе данных
	rows, err := r.db.Query(ctx,
		`SELECT 
            u.username, 
            p.age, 
            p.avatar_url, 
            p.about_me, 
			p.gender,
            p.looking_for, 
            p.created_at, 
            p.updated_at 
        FROM 
            profiles p
        JOIN 
            users u ON p.user_id = u.id
        WHERE 
            p.updated_at >= $1`,
		sevenDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*models.UserProfile

	// Итерируем по результатам запроса
	for rows.Next() {
		var username string
		var age int
		var avatarUrl string
		var aboutMe string
		var gender string
		var lookingFor string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&username, &age, &avatarUrl, &aboutMe, &gender, &lookingFor, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		profiles = append(profiles, &models.UserProfile{
			Username:   username,
			Age:        age,
			AvatarUrl:  avatarUrl,
			AboutMe:    aboutMe,
			Gender:     gender,
			LookingFor: lookingFor,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	// Проверяем, были ли ошибки при итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return profiles, nil
}

func (r *Repository) DeleteUserByID(ctx context.Context, userId int) error {

	_, err := r.db.Exec(ctx, `
    DELETE FROM users
    WHERE id = $1;`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	return nil
}

func (r *Repository) AddResponse(ctx context.Context, userId, profileId string) error {

	err := r.db.QueryRow(ctx,
		"INSERT INTO responses (profile_id, responder_id, status, created_At) VALUES ($1, $2, $3, $4)",
		profileId, userId, "ожидание", time.Now())

	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

func (r *Repository) ConfirmResponse(ctx context.Context, responseId string) error {

	_, err := r.db.Exec(ctx,
		"UPDATE responses SET status = $1 WHERE id = $1",
		"одобрено", responseId)

	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

func (r *Repository) RejectResponse(ctx context.Context, responseId string) error {

	_, err := r.db.Exec(ctx,
		"UPDATE responses SET status = $1 WHERE id = $1",
		"отказано", responseId)

	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

func (r *Repository) GetIncomingResponses(ctx context.Context, userId int) ([]*models.UserResponse, error) {

	rows, err := r.db.Query(ctx,
		`SELECT 
            r.status, 
			u.username
            p.age, 
            p.avatar_url, 
            p.about_me, 
			p.gender,
            p.looking_for, 
            p.created_at, 
            p.updated_at 
        FROM 
            responses r
        JOIN 
            profiles p ON p.id = r.responder_id
		JOIN 
            users u ON p.user_id = u.id
        WHERE 
            r.responder_id = $1`,
		userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}
	defer rows.Close()

	var incomingResponses []*models.UserResponse

	for rows.Next() {
		var status string
		var username string
		var age int
		var avatarUrl string
		var aboutMe string
		var gender string
		var lookingFor string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&status, &username, &age, &avatarUrl, &aboutMe, &gender, &lookingFor, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		incomingResponses = append(incomingResponses, &models.UserResponse{
			Status: status,
			Responder: &models.UserProfile{
				Username:   username,
				Age:        age,
				AvatarUrl:  avatarUrl,
				AboutMe:    aboutMe,
				Gender:     gender,
				LookingFor: lookingFor,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			},
		})
	}

	return incomingResponses, nil
}

func (r *Repository) GetOutgoingResponses(ctx context.Context, userId int) ([]*models.UserResponse, error) {

	rows, err := r.db.Query(ctx,
		`SELECT 
            r.status, 
			u.username
            p.age, 
            p.avatar_url, 
            p.about_me, 
			p.gender,
            p.looking_for, 
            p.created_at, 
            p.updated_at 
        FROM 
            responses r
        JOIN 
            profiles p ON p.id = r.responder_id
		JOIN 
            users u ON p.user_id = u.id
        WHERE 
            r.profile_id = $1`,
		userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}
	defer rows.Close()

	var outgoingResponses []*models.UserResponse

	for rows.Next() {
		var status string
		var username string
		var age int
		var avatarUrl string
		var aboutMe string
		var gender string
		var lookingFor string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&status, &username, &age, &avatarUrl, &aboutMe, &gender, &lookingFor, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		outgoingResponses = append(outgoingResponses, &models.UserResponse{
			Status: status,
			Responder: &models.UserProfile{
				Username:   username,
				Age:        age,
				AvatarUrl:  avatarUrl,
				AboutMe:    aboutMe,
				Gender:     gender,
				LookingFor: lookingFor,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			},
		})
	}

	return outgoingResponses, nil
}

func (r *Repository) GetNotifications(ctx context.Context, userId int) ([]*models.Notification, error) {

	rows, err := r.db.Query(ctx,
		"SELECT message, is_read, created_at, type from notifications WHERE user_id = $1",
		userId)

	if err != nil {
		return nil, fmt.Errorf("failed to create response: %w", err)
	}

	var notifications []*models.Notification

	for rows.Next() {
		var message string
		var isRead bool
		var createdAt time.Time
		var notificationType int

		err := rows.Scan(&message, &isRead, &createdAt, &notificationType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notifications: %w", err)
		}

		notifications = append(notifications, &models.Notification{
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
			Type:      models.ConvertToNotidy(notificationType),
		})
	}

	return notifications, nil
}

func (r *Repository) AddNotification(ctx context.Context, userId int, message string, notificationType models.NotificationType) error {

	err := r.db.QueryRow(ctx,
		"INSERT INTO notifications (user_id, message, is_read, created_at, type) VALUES ($1, $2, $3, $4)",
		userId, message, false, time.Now(), notificationType.ToInt())

	if err != nil {
		return fmt.Errorf("failed to add notification: %w", err)
	}

	return nil
}
