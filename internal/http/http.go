package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	database2 "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiberApp(db *database.ECampusDB) *fiber.App {
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) (err error) {
		query, _, _ := db.QB.From("users").Limit(10).ToSQL()
		var users []database2.User
		err = db.Conn.Select(&users, query)
		if err != nil {
			return
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": users,
		})
	})

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
