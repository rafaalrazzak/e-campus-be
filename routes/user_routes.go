package routes

import (
	"ecampus/internal/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(r *gin.Engine) {
	userHandler := handlers.UserHandler{}

	userGroup := r.Group("/users")
	{
		userGroup.GET("", userHandler.GetAllUsers)       // Get all users
		userGroup.POST("", userHandler.CreateUser)       // Create a new user
		userGroup.GET("/:id", userHandler.GetUser)       // Get a user by ID
		userGroup.PUT("/:id", userHandler.UpdateUser)    // Update a user by ID
		userGroup.DELETE("/:id", userHandler.DeleteUser) // Delete a user by ID
	}
}
