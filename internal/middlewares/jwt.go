package middlewares

import (
	"fmt"
	netHttp "net/http"
	"time"

	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthorization struct {
	db *database.AuthDB
}

func (m *JWTAuthorization) Auth(
	ctx *fiber.Ctx,
	scope httpModels.Scope,
) (int, error) {
	tokenString := ctx.Get("Authorization")
	if tokenString == "" {
		return netHttp.StatusUnauthorized, nil
	}

	claims := new(httpModels.Claims)
	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return config.ServerAuthSecret, nil
		},
	)
	if err != nil || !token.Valid || time.Now().UTC().After(claims.ExpiresAt.UTC()) {
		return netHttp.StatusUnauthorized, nil
	}

	if claims.Scope != scope {
		return netHttp.StatusForbidden, nil
	}

	exists, err := m.db.Exists(ctx.Context(), tokenString)
	if err != nil {
		return 0, nil
	}
	if !exists {
		return netHttp.StatusUnauthorized, nil
	}

	ctx.Locals("telegram_id", claims.TelegramID)

	return 0, nil
}

func (m *JWTAuthorization) Middleware(scope httpModels.Scope) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		code, err := m.Auth(ctx, scope)
		if err != nil {
			return err
		}
		if code != 0 {
			return ctx.SendStatus(code)
		}
		return ctx.Next()
	}
}

func NewJWTAuthorization(db *database.AuthDB) *JWTAuthorization {
	return &JWTAuthorization{db: db}
}
