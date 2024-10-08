package routes

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	domain "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	"github.com/rafaalrazzak/e-campus-be/internal/utils"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"net/http"
)

func AuthRoutes(db *database.ECampusDB, redisClient *redis.ECampusRedis) fiber.Router {
	authRoute := fiber.New().Group("/auth")

	authRoute.Post("/login", func(c *fiber.Ctx) error {
		var inputUser domain.User

		if err := c.BodyParser(&inputUser); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		if inputUser.Email == "" && inputUser.Username == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Email or username is required")
		}

		if inputUser.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Password is required")
		}

		userQuery := db.QB.From("users").Where(
			goqu.Or(
				goqu.Ex{"username": inputUser.Username},
				goqu.Ex{"email": inputUser.Email},
			),
		)

		var dbUser domain.User
		sql, _, _ := userQuery.ToSQL()

		if err := db.Conn.Get(&dbUser, sql); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		if !utils.VerifyData(dbUser.Password, inputUser.Password) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		// Generate session token
		sessionToken, err := utils.GenerateSessionToken()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate session token")
		}

		// Store session in Redis
		// TODO: Implement Redis session storage

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"token":   sessionToken,
		})
	})

	authRoute.Post("/register", func(c *fiber.Ctx) error {
		var newUser domain.User

		if err := c.BodyParser(&newUser); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		// Validate input
		if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Username, email, and password are required")
		}

		// Check if username or email already exists
		existsQuery := db.QB.From("users").Where(
			goqu.Or(
				goqu.Ex{"username": newUser.Username},
				goqu.Ex{"email": newUser.Email},
			),
		)

		var existingUser domain.User
		sql, _, _ := existsQuery.ToSQL()

		err := db.Conn.Get(&existingUser, sql)
		if err == nil {
			return fiber.NewError(fiber.StatusConflict, "Username or email already exists")
		}

		// Hash password
		newUser.Password = utils.HashData(newUser.Password)

		// Generate user ID
		newUser.ID = utils.GenerateId()

		// Insert new user
		insertQuery, _, _ := db.QB.Insert("users").Rows(newUser).ToSQL()

		_, err = db.Conn.Exec(insertQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"userId":  newUser.ID,
		})
	})

	return nil
}
