package controllers

import (
	netHttp "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/randomowo-dev/telegram-films-bot/internal/services"
)

type ConfigController struct {
	configService *services.ConfigService
}

func (c *ConfigController) Config(ctx *fiber.Ctx) error {
	data, err := c.configService.GetConfig(ctx.UserContext())
	if err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError)
	}

	return ctx.JSON(data)
}

func NewConfigController(configService *services.ConfigService) *ConfigController {
	return &ConfigController{configService: configService}
}
