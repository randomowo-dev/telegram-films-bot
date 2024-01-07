package controllers

import (
	netHttp "net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	"github.com/randomowo-dev/telegram-films-bot/internal/services"

	_ "embed"
)

type UserController struct {
	userService *services.UserService
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	offset, _ := strconv.ParseInt(ctx.Query("offset", "0"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit", "10"), 10, 64)
	if limit == 0 {
		limit = 10
	}

	users, haveNext, err := c.userService.ListUsers(ctx.UserContext(), offset, limit)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError)
	}

	return ctx.JSON(
		map[string]any{
			"users":     users,
			"have_next": haveNext,
		},
	)
}

func (c *UserController) UpdateUserRole(ctx *fiber.Ctx) error {
	telegramID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusBadRequest)
	}
	role := dbModels.ParseRole(ctx.Query("role"))

	if err = c.userService.UpdateUserRole(ctx.UserContext(), telegramID, role); err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError)
	}

	return ctx.SendStatus(netHttp.StatusOK)
}

func NewAdminController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}
