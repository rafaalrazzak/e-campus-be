package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
)

func SetupRoutes(app *fiber.App, db *database.ECampusDB, redisDB *redis.ECampusRedisDB, config config.Config) {
	SetupAuthRoutes(app, db, redisDB, config)
	SetupUserRoutes(app, db, redisDB, config)
}
