package routes

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	database2 "github.com/rafaalrazzak/e-campus-be/internal/domain/database"
	"github.com/rafaalrazzak/e-campus-be/internal/utils"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"net/http"
)

func checkUserExists(db *database.ECampusDB, field, value string) (bool, error) {
	query, _, _ := db.QB.From("users").Where(goqu.Ex{field: value}).ToSQL()
	var existingUser database2.User
	err := db.Conn.Get(&existingUser, query)
	if err == nil {
		return true, nil
	}
	return false, err
}

func UserRoutes(app *fiber.App, db *database.ECampusDB) fiber.Router {
	userRoute := app.Group("/users")

	userRoute.Get("/", func(c *fiber.Ctx) error {
		query, _, _ := db.QB.From("users").Limit(10).ToSQL()
		var users []database2.User
		err := db.Conn.Select(&users, query)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": users,
		})
	})

	userRoute.Post("/", func(c *fiber.Ctx) error {
		var user database2.User
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		user.ID = utils.GenerateId()
		user.Password = utils.HashData(user.Password)

		tx, err := db.Conn.Begin()
		if err != nil {
			return err
		}
		defer func(tx *sql.Tx) {
			_ = tx.Rollback()
		}(tx)

		if exists, _ := checkUserExists(db, "username", user.Username); exists {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"message": "username already exists",
			})
		}

		if exists, _ := checkUserExists(db, "email", user.Email); exists {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"message": "email already exists",
			})
		}

		query, _, _ := db.QB.Insert("users").Rows(user).Returning("id").ToSQL()
		err = tx.QueryRow(query).Scan(&user.ID)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		return c.SendStatus(http.StatusCreated)
	})

	userRoute.Patch("/:id", func(c *fiber.Ctx) error {
		var user database2.User
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		tx, err := db.Conn.Begin()
		if err != nil {
			return err
		}
		defer func(tx *sql.Tx) {
			_ = tx.Rollback()
		}(tx)

		if exists, _ := checkUserExists(db, "id", c.Params("id")); !exists {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "user not found",
			})
		}

		// Check if username, email already exists
		if user.Username != "" {
			if exists, _ := checkUserExists(db, "username", user.Username); exists {
				return c.Status(http.StatusConflict).JSON(fiber.Map{
					"message": "username already exists",
				})
			}
		}

		if user.Email != "" {
			if exists, _ := checkUserExists(db, "email", user.Email); exists {
				return c.Status(http.StatusConflict).JSON(fiber.Map{
					"message": "email already exists",
				})
			}
		}

		updateFields := goqu.Record{}
		if user.Username != "" {
			updateFields["username"] = user.Username
		}
		if user.Email != "" {
			updateFields["email"] = user.Email
		}
		if user.Password != "" {
			updateFields["password"] = utils.HashData(user.Password)
		}

		query, _, _ := db.QB.Update("users").Set(updateFields).Where(goqu.Ex{"id": c.Params("id")}).ToSQL()
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})
	return nil
}
