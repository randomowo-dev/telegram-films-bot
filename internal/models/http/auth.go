package http

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	TelegramID int64 `json:"telegram_id"`
	Scope      Scope `json:"scope"`
	jwt.StandardClaims
}

type Scope uint8

const (
	Refresh Scope = iota
	Api
)

type AuthUser struct {
	TelegramID int64
	Username   string
	AuthDate   time.Time
}

type NewTokenResponse struct {
	Token                  string `json:"token"`
	TokenExpiration        int64  `json:"token_expiration"`
	RefreshToken           string `json:"refresh_token"`
	RefreshTokenExpiration int64  `json:"refresh_token_expiration"`
}

type RefreshToken struct {
	TelegramID int64
	Token      string
}
