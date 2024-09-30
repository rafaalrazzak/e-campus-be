package services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// jsonResponse writes a JSON response with a status code
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func handleError(c *gin.Context, err any, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err})
}
