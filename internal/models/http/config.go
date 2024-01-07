package http

import (
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
)

type Config struct {
	AvailableRoles []dbModels.Role `json:"available_roles"`
}
