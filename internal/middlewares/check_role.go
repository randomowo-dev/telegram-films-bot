package middlewares

import (
	netHttp "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
)

type RoleChecker struct {
	userDB *database.UserDB
}

func (c *RoleChecker) Check(ctx *fiber.Ctx, role dbModels.Role) (bool, error) {
	telegramID, _ := ctx.Locals("telegram_id").(int64)

	user, err := c.userDB.GetByTelegramID(ctx.UserContext(), telegramID)
	if err != nil {
		return false, err
	}

	return user != nil && user.Role >= role, nil
}

func (c *RoleChecker) Middleware(role dbModels.Role) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ok, err := c.Check(ctx, role)
		if err != nil {
			return ctx.SendStatus(netHttp.StatusInternalServerError)
		}

		if !ok {
			return ctx.SendStatus(netHttp.StatusForbidden)
		}

		return ctx.Next()
	}
}

func NewRoleChecker(userDB *database.UserDB) *RoleChecker {
	return &RoleChecker{userDB: userDB}
}
