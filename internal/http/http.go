package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"go.uber.org/fx"
)

func NewFiberApp() *fiber.App {
	app := fiber.New()

	return app
}

func ServeHTTP(lc fx.Lifecycle, app *fiber.App, config config.Config) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return app.Listen(fmt.Sprintf(":%s", config.ServerPort))
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
}
