package services

import (
	"context"

	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
)

type ConfigService struct{}

func (s *ConfigService) GetConfig(ctx context.Context) (*httpModels.Config, error) {
	return &httpModels.Config{
		AvailableRoles: dbModels.AvailableRoles,
	}, nil
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}
