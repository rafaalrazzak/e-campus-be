package middleware

import (
	"github.com/rafaalrazzak/e-campus-be/internal/services"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
)

// Config holds all middleware dependencies
type Config struct {
	DB     *database.ECampusDB
	Config config.Config
}

// AuthorizationMiddleware validates user authorization based on token
func AuthorizationMiddleware(db *database.ECampusDB, redisDB *redis.ECampusRedisDB, config config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authService := services.NewAuthService(db, redisDB, config)
		userService := services.NewUserService(db)

		// Get token from Authorization header
		token := extractToken(c)

		if token == "" {
			return fiber.NewError(http.StatusUnauthorized, "Missing or invalid authorization token")
		}

		userId, err := authService.GetSession(token)

		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, "Invalid or expired token")
		}

		// Retrieve user data
		userData, err := userService.GetUserByID(*userId)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		// Store user data and claims in context
		c.Locals("userData", userData)
		return c.Next()
	}
}

// RoleAuthMiddleware checks if the user has the required role
func RoleAuthMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData, ok := c.Locals("userData").(map[string]interface{})
		if !ok {
			return fiber.NewError(http.StatusUnauthorized, "User data not found")
		}

		userRole, ok := userData["role"].(string)
		if !ok {
			return fiber.NewError(http.StatusForbidden, "User role not found")
		}

		for _, role := range requiredRoles {
			if role == userRole {
				return c.Next()
			}
		}

		return fiber.NewError(http.StatusForbidden, "Insufficient permissions")
	}
}

// RateLimitMiddleware implements a basic rate limiting
func RateLimitMiddleware(requests int, duration time.Duration) fiber.Handler {
	// Simple in-memory store for rate limiting
	type client struct {
		count    int
		lastSeen time.Time
	}
	clients := make(map[string]*client)

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		if clients[ip] == nil {
			clients[ip] = &client{count: 1, lastSeen: now}
			return c.Next()
		}

		if now.Sub(clients[ip].lastSeen) > duration {
			clients[ip].count = 1
			clients[ip].lastSeen = now
			return c.Next()
		}

		if clients[ip].count >= requests {
			return fiber.NewError(http.StatusTooManyRequests, "Rate limit exceeded")
		}

		clients[ip].count++
		return c.Next()
	}
}

// extractToken helper function to extract token from Authorization header
func extractToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	if token == "" {
		return ""
	}

	return token
}
