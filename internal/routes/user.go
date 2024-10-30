// internal/routes/user_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/controllers"
	"github.com/rafaalrazzak/e-campus-be/internal/middleware"
	"github.com/rafaalrazzak/e-campus-be/internal/services"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
)

func SetupUserRoutes(router fiber.Router, db *database.ECampusDB, redisDB *redis.ECampusRedisDB, config config.Config) {
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	users := router.Group("/users")

	// Public routes
	users.Get("/", userController.GetUsers())

	me := users.Group("/me")
	me.Use(middleware.AuthorizationMiddleware(db, redisDB, config))
	me.Get("/", userController.GetCurrentUser())

	// Protected routes
	users.Post("/", middleware.RoleAuthMiddleware("admin"), userController.CreateUser())
	users.Put("/:id", middleware.RoleAuthMiddleware("admin"), userController.UpdateUser())
	users.Delete("/:id", middleware.RoleAuthMiddleware("admin"), userController.DeleteUser())
}
