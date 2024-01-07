package services

import (
	"context"
	"time"

	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userDB *database.UserDB
	authDB *database.AuthDB
}

func (s *AuthService) generateTokens(telegramID int64) (*httpModels.NewTokenResponse, error) {
	var err error
	auth := new(httpModels.NewTokenResponse)

	authClaims := &httpModels.Claims{
		Scope:      httpModels.Api,
		TelegramID: telegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
	}
	auth.TokenExpiration = authClaims.ExpiresAt.UnixNano()

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	auth.Token, err = authToken.SignedString(config.ServerAuthSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := &httpModels.Claims{
		Scope:      httpModels.Auth,
		TelegramID: telegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		},
	}
	auth.RefreshTokenExpiration = refreshClaims.ExpiresAt.UnixNano()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	auth.RefreshToken, err = refreshToken.SignedString(config.ServerAuthSecret)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *AuthService) AuthUser(ctx context.Context, user *httpModels.AuthUser) (*httpModels.NewTokenResponse, error) {
	authData, err := s.generateTokens(user.TelegramID)
	if err != nil {
		return nil, err
	}

	userID, err := s.userDB.UpdateByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return nil, err
	}

	if userID == "" {
		dbUser := &dbModels.User{
			TelegramID: user.TelegramID,
			Username:   user.Username,
			LastAuth:   user.AuthDate,
		}

		if err = s.userDB.Add(ctx, dbUser); err != nil {
			return nil, err
		}

		userID = dbUser.ID
	}

	auth := &dbModels.Auth{
		UserID:       userID,
		Token:        authData.Token,
		RefreshToken: authData.RefreshToken,
	}
	if err := s.authDB.Add(ctx, auth); err != nil {
		// FIXME: do something with updated data above with error here
		return nil, err
	}

	return &httpModels.NewTokenResponse{
		Token:                  authData.Token,
		TokenExpiration:        authData.TokenExpiration,
		RefreshToken:           authData.RefreshToken,
		RefreshTokenExpiration: authData.RefreshTokenExpiration,
	}, nil
}

func (s *AuthService) RefreshUserToken(
	ctx context.Context,
	telegramID int64,
) (*httpModels.NewTokenResponse, error) {
	user, err := s.userDB.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	_ = s.authDB.DeleteByUserID(ctx, user.ID)

	return s.AuthUser(
		ctx, &httpModels.AuthUser{
			TelegramID: user.TelegramID,
			Username:   user.Username,
			AuthDate:   user.LastAuth,
		},
	)
}

func (s *AuthService) LogOut(
	ctx context.Context,
	telegramID int64,
) error {
	user, err := s.userDB.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	if err := s.authDB.DeleteByUserID(ctx, user.ID); err != nil {
		return err
	}

	return nil
}

func NewAuthService(userDB *database.UserDB, authDB *database.AuthDB) *AuthService {
	return &AuthService{
		userDB: userDB,
		authDB: authDB,
	}
}
