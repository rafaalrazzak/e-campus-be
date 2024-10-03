package ecampus

import (
	"context"

	"ecampus-be/bunapp"
)

func init() {
	bunapp.OnStart("ecampus.init", func(ctx context.Context, app *bunapp.App) error {
		router := app.Router()

		welcomeHandler := NewWelcomeHandler(app)
		userHandler := NewUserHandler(app)
		authHandler := NewAuthHandler(app)

		router.GET("/", welcomeHandler.Welcome)
		router.GET("/hello", welcomeHandler.Hello)

		users := router.NewGroup("/users")
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)

		auth := router.NewGroup("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)

		return nil
	})
}
