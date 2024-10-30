package controllers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/rafaalrazzak/e-campus-be/internal/domain/models"
	"github.com/rafaalrazzak/e-campus-be/internal/services"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		limit := ctx.QueryInt("limit", 10)
		page := ctx.QueryInt("page", 1)
		offset := (page - 1) * limit

		// Parse filters from query parameters
		filters := make(map[string]interface{})
		if role := ctx.Query("role"); role != "" {
			filters["role"] = role
		}
		if status := ctx.Query("status"); status != "" {
			filters["status"] = status
		}
		if dept := ctx.Query("department_code"); dept != "" {
			filters["department_code"] = dept
		}

		params := services.UserFilters{
			Limit:   limit,
			Offset:  offset,
			Filters: filters,
		}

		response, err := c.userService.GetUsers(params)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch users")
		}

		return ctx.JSON(response)
	}
}

func (c *UserController) GetCurrentUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("userData")

		return ctx.JSON(user)
	}
}

func (c *UserController) CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user models.User
		if err := ctx.BodyParser(&user); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		if err := c.userService.CreateUser(&user); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return ctx.Status(http.StatusCreated).JSON(user)
	}
}

func (c *UserController) UpdateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Params("id")
		if userID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
		}

		updates := make(map[string]interface{})
		if err := ctx.BodyParser(&updates); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		if err := c.userService.UpdateUser(userID, updates); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "User not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func (c *UserController) DeleteUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Params("id")
		if userID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
		}

		if err := c.userService.DeleteUser(userID); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "User not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
		}

		return ctx.SendStatus(http.StatusNoContent)
	}
}
