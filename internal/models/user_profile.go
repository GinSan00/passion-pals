package model

import "time"

type UserProfile struct {
	Username   string    `json:"username"`
	Age        int       `json:"age"`
	AvatarUrl  string    `json:"avatar_url"`
	AboutMe    string    `json:"about_me"`
	Gender     string    `json:"gender"`
	LookingFor string    `json:"looking_for"`
	CreatedAt  time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt  time.Time `json:"updated_at" time_format:"2006-01-02T15:04:05Z07:00"`
}
