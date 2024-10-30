package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/controllers"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
)

// SetupAuthRoutes configures all authentication-related routes
func SetupAuthRoutes(app *fiber.App, db *database.ECampusDB, redisDB *redis.ECampusRedisDB, config config.Config) {
	authController := controllers.NewAuthController(db, redisDB, config)

	auth := app.Group("/auth")

	auth.Post("/login", authController.Login())
	//auth.Post("/forgot-password", controllers.ForgotPassword(db))
	//auth.Post("/reset-password", controllers.ResetPassword(db))
}
