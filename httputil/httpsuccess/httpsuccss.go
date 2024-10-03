package httpsuccess

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents the structure of a standard success response.
type SuccessResponse struct {
	Status  int         `json:"status"`  // HTTP status code
	Message string      `json:"message"` // A human-readable message
	Data    interface{} `json:"data"`    // Optional: any additional data
}

// New creates a new success response and writes it to the http.ResponseWriter.
func New(w http.ResponseWriter, status int, message string, data interface{}) {
	response := SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

//------------------------------------------------------------------------------

// Helper functions for common success responses.

func Created(w http.ResponseWriter, message string, data interface{}) error {
	New(w, http.StatusCreated, message, data)
	return nil
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	New(w, http.StatusOK, message, data)
}

func NoContent(w http.ResponseWriter, message string) error {
	New(w, http.StatusNoContent, message, nil)
	return nil
}
