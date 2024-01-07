package http

import (
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/pkg/transport/http/server"

	"github.com/gofiber/fiber/v2"
)

type AppServer struct {
	*server.Server
}

func (s *AppServer) Listen() error {
	return s.Server.Listen(config.ServerAddr)
}

func NewAppServer() *AppServer {
	return &AppServer{
		Server: server.NewServer(
			fiber.Config{
				AppName: config.AppName,
			},
		),
	}
}
