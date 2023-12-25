package middlewares

import (
	netHttp "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
)

func BasicAuthorization(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") != config.ServerBasicAuthToken {
		return ctx.SendStatus(netHttp.StatusUnauthorized)
	}

	return ctx.Next()
}
