// internal/controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/domain/models"
	"github.com/rafaalrazzak/e-campus-be/internal/services"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(db *database.ECampusDB, redisClient *redis.ECampusRedisDB, cfg config.Config) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(db, redisClient, cfg),
	}
}

// Login handles user authentication
func (ctrl *AuthController) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var inputUser models.User

		if err := c.BodyParser(&inputUser); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		token, err := ctrl.authService.AuthenticateUser(inputUser)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}
}
