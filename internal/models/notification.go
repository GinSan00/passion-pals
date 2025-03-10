package model

import "time"

type Notification struct {
	Message   string
	IsRead    bool
	CreatedAt time.Time
	Type      NotificationType
}
