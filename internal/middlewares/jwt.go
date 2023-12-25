package middlewares

import (
	"fmt"
	netHttp "net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
)

func JWTAuthorizationMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	if tokenString == "" {
		return ctx.SendStatus(netHttp.StatusUnauthorized)
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
	if err != nil || !token.Valid || time.Now().UTC().After(time.Unix(0, claims.ExpiresAt)) {
		return ctx.SendStatus(netHttp.StatusUnauthorized)
	}

	if claims.Scope != httpModels.Api {
		return ctx.SendStatus(netHttp.StatusForbidden)
	}

	ctx.Locals("telegram_id", claims.TelegramID)

	return ctx.Next()
}
