package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/routes"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp(db *database.ECampusDB) *fiber.App {
	app := fiber.New()

	routes.UserRoutes(app, db)

	return app
}

func ServeHTTP(lc fx.Lifecycle, app *fiber.App, config config.Config) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				go func() {
					err = app.Listen(fmt.Sprintf(":%s", config.ServerPort))
					if err != nil {
						zap.L().Fatal("failed to start http server:", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
}
