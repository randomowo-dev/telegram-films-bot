package http

import (
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
)

type UpdateUser struct {
	UserID string        `json:"user_id"`
	Role   dbModels.Role `json:"role"`
}

type User struct {
	TelegramID int64         `json:"telegram_id"`
	Username   string        `json:"username"`
	Role       dbModels.Role `json:"role"`
}
