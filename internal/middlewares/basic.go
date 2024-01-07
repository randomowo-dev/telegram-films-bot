package middlewares

import (
	netHttp "net/http"

	"github.com/randomowo-dev/telegram-films-bot/internal/config"

	"github.com/gofiber/fiber/v2"
)

func BasicAuthorization(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") != config.ServerBasicAuthToken {
		return ctx.SendStatus(netHttp.StatusUnauthorized)
	}

	return ctx.Next()
}
