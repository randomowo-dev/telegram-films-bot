package server

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Middleware(path string, middleware fiber.Handler) {
	s.app.Use(path, middleware)
}

func (s *Server) AddRoute(method, path string, handler fiber.Handler) {
	s.app.Add(method, path, handler)
}

func (s *Server) Group(prefix string, middlewares ...fiber.Handler) fiber.Router {
	return s.app.Group(prefix, middlewares...)
}

func NewServer(config fiber.Config) *Server {
	return &Server{
		app: fiber.New(config),
	}
}
