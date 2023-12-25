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

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
	"github.com/randomowo-dev/telegram-films-bot/internal/services"
)

type AuthController struct {
	authService *services.AuthService
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

	if claims.Scope != httpModels.Refresh {
		return ctx.SendStatus(netHttp.StatusForbidden)
	}

	data, err := c.authService.RefreshUserToken(ctx.UserContext(), claims.TelegramID)
	if err != nil {
		return ctx.SendStatus(netHttp.StatusInternalServerError) // FIXME
	}

	return ctx.JSON(data)
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}
