package ecampus

import (
	"context"

	"github.com/go-bun/bun-starter-kit/bunapp"
)

func init() {
	bunapp.OnStart("ecampus.init", func(ctx context.Context, app *bunapp.App) error {
		router := app.Router()

		welcomeHandler := NewWelcomeHandler(app)
		userHandler := NewUserHandler(app)

		router.GET("/", welcomeHandler.Welcome)
		router.GET("/hello", welcomeHandler.Hello)

		users := router.NewGroup("/users")
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.Get)
		users.POST("", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)

		return nil
	})
}
