package users

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	_ "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	database2 "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	"github.com/rafaalrazzak/e-campus-be/internal/utils"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"net/http"
)

func GetUsers(db *database.ECampusDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []database2.User
		query, _, _ := db.QB.From("users").ToSQL()
		if err := db.Conn.Select(&users, query); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get users")
		}

		return c.JSON(users)
	}
}

func CreateUser(db *database.ECampusDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user database2.User
		if err := c.BodyParser(&user); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		user.ID = utils.GenerateId()
		passwordHash, err := utils.HashData(user.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
		}
		user.Password = passwordHash

		exists, err := checkUserExists(db, "email", user.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to check user existence")
		}
		if exists {
			return fiber.NewError(fiber.StatusConflict, "Email already exists")
		}

		query, _, _ := db.QB.Insert("users").Rows(user).Returning("id").ToSQL()
		if err := db.Conn.QueryRow(query).Scan(&user.ID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(http.StatusCreated)
	}
}

func UpdateUser(db *database.ECampusDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user database2.User
		if err := c.BodyParser(&user); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		exists, err := checkUserExists(db, "id", c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to check user existence")
		}
		if !exists {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		if user.Email != "" {
			exists, err := checkUserExists(db, "email", user.Email)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to check email existence")
			}
			if exists {
				return fiber.NewError(fiber.StatusConflict, "Email already exists")
			}
		}

		updateFields := goqu.Record{}
		if user.Email != "" {
			updateFields["email"] = user.Email
		}
		if user.Password != "" {
			passwordHash, err := utils.HashData(user.Password)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
			}
			updateFields["password"] = passwordHash
		}

		query, _, _ := db.QB.Update("users").Set(updateFields).Where(goqu.Ex{"id": c.Params("id")}).ToSQL()
		if _, err := db.Conn.Exec(query); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
		}

		return c.SendStatus(http.StatusNoContent)
	}
}

func checkUserExists(db *database.ECampusDB, field, value string) (bool, error) {
	query, _, _ := db.QB.From("users").Where(goqu.Ex{field: value}).ToSQL()
	var existingUser database2.User
	err := db.Conn.Get(&existingUser, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No user found
		}
		return false, err // Some other error occurred
	}
	return true, nil // User exists
}
