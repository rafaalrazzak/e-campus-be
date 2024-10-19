package auth

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/constants"
	database2 "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	"github.com/rafaalrazzak/e-campus-be/internal/utils"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"net/http"
)

func Login(db *database.ECampusDB, redisClient *redis.ECampusRedisDB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var inputUser database2.User

		if err := c.BodyParser(&inputUser); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		if err := validateLoginInput(inputUser); err != nil {
			return err
		}

		dbUser, err := getUserFromDB(db, inputUser.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		if err := verifyPassword(dbUser.Password, inputUser.Password); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		sessionId := utils.GenerateSessionToken()
		sessionToken := fmt.Sprintf("%d::%s", dbUser.ID, sessionId)

		token, err := utils.GenerateSessionEncryption(sessionToken, cfg)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		ctx := context.Background()

		userData := map[string]interface{}{
			"id":    dbUser.ID,
			"email": dbUser.Email,
			"name":  dbUser.Name,
			"group": dbUser.Group,
			"role":  dbUser.Role,
			"major": dbUser.Major,
			"year":  dbUser.Year,
		}

		redisKey := fmt.Sprintf(constants.Redis.SessionKey, dbUser.ID, sessionId)
		if err := redisClient.Client.HSet(ctx, redisKey, userData).Err(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to store session data")
		}
		if err := redisClient.Client.Expire(ctx, redisKey, constants.App.SessionExpiration).Err(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to set session expiration")
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}
}

func validateLoginInput(user database2.User) error {
	if user.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}
	if user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Password is required")
	}
	return nil
}

func getUserFromDB(db *database.ECampusDB, email string) (database2.User, error) {
	var dbUser database2.User
	query := db.QB.From("users").Where(goqu.Ex{"email": email})
	sql, _, _ := query.ToSQL()
	err := db.Conn.Get(&dbUser, sql)
	return dbUser, err
}

func verifyPassword(hashedPassword, inputPassword string) error {
	isValid, err := utils.VerifyData(hashedPassword, inputPassword)
	if err != nil {
		return err
	}
	if !isValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}
	return nil
}
