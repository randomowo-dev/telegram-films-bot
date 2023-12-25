package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
)

type AuthService struct {
	userDB *database.UserDB
}

func (s *AuthService) generateTokens(telegramID int64) (*httpModels.NewTokenResponse, error) {
	var err error
	auth := new(httpModels.NewTokenResponse)

	auth.TokenExpiration = time.Now().Add(time.Hour).UnixNano()
	authClaims := &httpModels.Claims{
		Scope:      httpModels.Api,
		TelegramID: telegramID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: auth.TokenExpiration,
		},
	}
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	auth.Token, err = authToken.SignedString(config.ServerAuthSecret)
	if err != nil {
		return nil, err
	}

	auth.RefreshTokenExpiration = time.Now().Add(24 * time.Hour).UnixNano()
	refreshClaims := &httpModels.Claims{
		Scope:      httpModels.Refresh,
		TelegramID: telegramID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: auth.RefreshTokenExpiration,
		},
	}
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

	dbUser := &dbModels.User{
		TelegramID: user.TelegramID,
		Username:   user.Username,
		LastAuth:   user.AuthDate,
	}

	if err := s.userDB.UpdateByTelegramID(ctx, dbUser); err != nil {
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
	user, err := s.userDB.FindByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	return s.AuthUser(
		ctx, &httpModels.AuthUser{
			TelegramID: user.TelegramID,
			Username:   user.Username,
			AuthDate:   user.LastAuth,
		},
	)
}

func NewAuthService(userDB *database.UserDB) *AuthService {
	return &AuthService{
		userDB: userDB,
	}
}
