package response

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Message string `json:"message"`
	// Data    interface{} `json:"data,omitempty"`
	Data *fiber.Map `json.data`
	// Error  error       `json:"-"`
	Status int `json:"status"`
}
