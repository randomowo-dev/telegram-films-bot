package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	netHttp "net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/internal/middlewares"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
	"github.com/randomowo-dev/telegram-films-bot/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService    *services.AuthService
	authMiddleware *middlewares.JWTAuthorization
}

func (c *AuthController) AuthUser(ctx *fiber.Ctx) error {
	var err error
	dataString := make([]string, 0)
	for k, v := range ctx.Queries() {
		if k == "hash" {
			continue
		}
		dataString = append(dataString, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(dataString)

	h := hmac.New(sha256.New, config.TelegramBotToken.Sum(nil))
	h.Write([]byte(strings.Join(dataString, "\n")))
	userHash := hex.EncodeToString(h.Sum(nil))

	queryHash := ctx.Query("hash")
	if userHash != queryHash {
		return ctx.SendStatus(netHttp.StatusBadRequest)
	}

	user := new(httpModels.AuthUser)
	user.AuthDate = time.Unix(int64(ctx.QueryInt("auth_date")), 0)
	if time.Now().UTC().Add(-24 * time.Hour).After(user.AuthDate) { // auth date is time when log in button pressed
		return ctx.SendStatus(netHttp.StatusBadRequest)
	}

	user.TelegramID, err = strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusBadRequest)
	}
	user.Username = ctx.Query("username")

	data, err := c.authService.AuthUser(ctx.UserContext(), user)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError) // FIXME
	}

	return ctx.JSON(data)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	code, err := c.authMiddleware.Auth(ctx, httpModels.AuthScope)
	if err != nil {
		return err
	}
	if code != 0 {
		return ctx.SendStatus(code)
	}

	telegramID, _ := ctx.Locals("telegram_id").(int64)
	data, err := c.authService.RefreshUserToken(ctx.UserContext(), telegramID)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError) // FIXME
	}

	return ctx.JSON(data)
}

func (c *AuthController) LogOut(ctx *fiber.Ctx) error {
	code, err := c.authMiddleware.Auth(ctx, httpModels.AuthScope)
	if err != nil {
		return err
	}
	if code != 0 {
		return ctx.SendStatus(code)
	}

	telegramID, _ := ctx.Locals("telegram_id").(int64)
	if err := c.authService.LogOut(ctx.UserContext(), telegramID); err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError)
	}

	return ctx.SendStatus(netHttp.StatusOK)
}

func NewAuthController(
	authService *services.AuthService,
	authMiddleware *middlewares.JWTAuthorization,
) *AuthController {
	return &AuthController{
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}
