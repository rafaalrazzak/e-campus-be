// http/http.go

package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rafaalrazzak/e-campus-be/internal/routes"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp(db *database.ECampusDB, redisClient *redis.ECampusRedisDB, cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	routes.SetupRoutes(app, db, redisClient, cfg)

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

func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
