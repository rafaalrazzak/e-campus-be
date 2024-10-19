package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/controllers/auth"
	"github.com/rafaalrazzak/e-campus-be/internal/controllers/users"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp(db *database.ECampusDB, redisClient *redis.ECampusRedisDB, cfg config.Config) *fiber.App {
	app := fiber.New()

	usersRoute := app.Group("/users")
	usersRoute.Get("/", users.GetUsers(db))
	usersRoute.Post("/", users.CreateUser(db))
	usersRoute.Patch("/:id", users.UpdateUser(db))

	authRoute := app.Group("/auth")
	authRoute.Post("/login", auth.Login(db, redisClient, cfg))

	return app
}

func ServeHTTP(lc fx.Lifecycle, app *fiber.App, config config.Config) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := app.Listen(fmt.Sprintf(":%s", config.ServerPort)); err != nil {
						zap.L().Fatal("Failed to start HTTP server:", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				// Gracefully shutdown Fiber app
				return app.Shutdown()
			},
		})
}
