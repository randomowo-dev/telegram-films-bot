package app

import (
	"context"
	netHttp "net/http"
	"sync"

	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/internal/controllers"
	"github.com/randomowo-dev/telegram-films-bot/internal/middlewares"
	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
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

	globalContext, cancelGlobalContextFunc := context.WithCancelCause(context.Background())

	db := nosql.NewDB()
	if err := db.Connect(globalContext); err != nil {
		panic(err)
	}

	userDB := database.NewUserDB(db)
	authDB := database.NewAuthDB(db)

	jwtAuth := middlewares.NewJWTAuthorization(authDB)
	roleChecker := middlewares.NewRoleChecker(userDB)

	manageGroup := server.Group(
		"/manage", basicauth.New(
			basicauth.Config{
				Users: map[string]string{
					config.ServerBasicUsername: config.ServerBasicPassword,
				},
			},
		),
	)
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

	// /api/v1/auth
	authController := controllers.NewAuthController(services.NewAuthService(userDB, authDB), jwtAuth)
	authGroup := groupV1.Group("/auth")
	authGroup.Add(netHttp.MethodGet, "/", authController.AuthUser)
	authGroup.Add(netHttp.MethodPut, "/refresh", authController.RefreshToken)
	authGroup.Add(netHttp.MethodPost, "/logout", authController.LogOut)

	// /api/v1/config
	configController := controllers.NewConfigController(services.NewConfigService())
	configGroup := groupV1.Group(
		"/config",
		jwtAuth.Middleware(httpModels.ApiScope),
		roleChecker.Middleware(dbModels.AdminRole),
	)
	configGroup.Add(netHttp.MethodGet, "/", configController.Config)

	// /api/v1/user
	userController := controllers.NewAdminController(services.NewUserService(userDB))
	userGroup := groupV1.Group(
		"/user",
		jwtAuth.Middleware(httpModels.ApiScope),
		roleChecker.Middleware(dbModels.AdminRole),
	)
	userGroup.Add(netHttp.MethodGet, "/", userController.List)
	userGroup.Add(netHttp.MethodPut, "/:id", userController.UpdateUserRole)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(globalContext context.Context, cancelContextFunc context.CancelCauseFunc, wg *sync.WaitGroup) {
		defer wg.Done()
		cancelContextFunc(server.Listen())
	}(globalContext, cancelGlobalContextFunc, wg)

	wg.Wait()
}
