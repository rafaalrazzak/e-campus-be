package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/route"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp() *fiber.App {
	app := fiber.New()

	return app
}

func RegisterRoutes(routes []route.Route, app *fiber.App) error {
	for _, r := range routes {
		app.Use(r)
	}

	return nil
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
