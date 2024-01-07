package app

import (
	"context"
	netHttp "net/http"

	"github.com/randomowo-dev/telegram-films-bot/internal/controllers"
	"github.com/randomowo-dev/telegram-films-bot/internal/middlewares"
	httpModels "github.com/randomowo-dev/telegram-films-bot/internal/models/http"
	"github.com/randomowo-dev/telegram-films-bot/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/randomowo-dev/telegram-films-bot/internal/database"
	"github.com/randomowo-dev/telegram-films-bot/internal/transport/http"
	"github.com/randomowo-dev/telegram-films-bot/internal/transport/nosql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
)

func Run() {
	reg := prometheus.NewRegistry()
	server := http.NewAppServer()
	// client := httpClient.NewClient()
	// kinopoiskApiUnofficialClient := http.NewKinopoiskApiUnofficialClient(client)
	globalContext, _ := context.WithCancel(context.Background())
	db := nosql.NewDB()
	if err := db.Connect(globalContext); err != nil {
		panic(err)
	}

	userDB := database.NewUserDB(db)
	authDB := database.NewAuthDB(db)
	// listDB := database.NewListDB(db)

	jwtAuth := middlewares.NewJWTAuthorization(authDB)
	authController := controllers.NewAuthController(services.NewAuthService(userDB, authDB), jwtAuth)

	manageGroup := server.Group("/manage", middlewares.BasicAuthorization)
	manageGroup.Add(
		netHttp.MethodGet, "/health", func(ctx *fiber.Ctx) error {
			ctx.Status(netHttp.StatusOK)
			return nil
		},
	)
	manageGroup.Add(
		netHttp.MethodGet,
		"/metrics",
		adaptor.HTTPHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})),
	)

	apiGroup := server.Group("/api")
	groupV1 := apiGroup.Group("/v1")

	authGroup := groupV1.Group("/auth")

	authGroup.Add(netHttp.MethodGet, "/", authController.AuthUser)
	authGroup.Add(netHttp.MethodPut, "/refresh", authController.RefreshToken)
	authGroup.Add(netHttp.MethodPost, "/logout", authController.LogOut)

	listGroup := groupV1.Group("list", jwtAuth.Middleware(httpModels.Api))

	listGroup.Add(
		netHttp.MethodGet, "/", func(ctx *fiber.Ctx) error {
			return ctx.SendString("test")
		},
	)

	_ = server.Listen()

	// listGroup.Add(
	// 	netHttp.MethodGet, "/", func(ctx *fiber.Ctx) error {
	// 		// TODO: get all
	// 		return nil
	// 	},
	// )
	// listGroup.Add(
	// 	netHttp.MethodGet, "/:id", func(ctx *fiber.Ctx) error {
	// 		// TODO: get by id
	// 		id := ctx.Params("id")
	// 		return nil
	// 	},
	// )
	// listGroup.Add(
	// 	netHttp.MethodPost, "/", func(ctx *fiber.Ctx) error {
	// 		// TODO: add new
	// 		return nil
	// 	},
	// )
	// listGroup.Add(
	// 	netHttp.MethodPut, "/:id", func(ctx *fiber.Ctx) error {
	// 		// TODO: update by id
	// 		id := ctx.Params("id")
	// 		return nil
	// 	},
	// )
	// listGroup.Add(
	// 	netHttp.MethodDelete, "/:id", func(ctx *fiber.Ctx) error {
	// 		// TODO: delete by id
	// 		id := ctx.Params("id")
	// 		return nil
	// 	},
	// )
}
