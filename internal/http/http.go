package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp() *fiber.App {
	app := fiber.New()

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
