package services

import (
	"context"

	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
)

type UserService struct {
	db *database.UserDB
}

func (s *UserService) ListUsers(ctx context.Context, offset int64, limit int64) ([]httpModels.User, bool, error) {
	users, haveNext, err := s.db.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, false, err
	}

	clearUsers := make([]httpModels.User, len(users))
	for index := range users {
		clearUsers[index] = httpModels.User{
			TelegramID: users[index].TelegramID,
			Username:   users[index].Username,
			Role:       users[index].Role,
		}
	}

	return clearUsers, haveNext, nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, telegramID int64, role dbModels.Role) error {
	dbUser := &dbModels.User{
		Role: role,
	}

	return s.db.UpdateByTelegramID(ctx, telegramID, dbUser)
}

func NewUserService(db *database.UserDB) *UserService {
	return &UserService{db: db}
}
