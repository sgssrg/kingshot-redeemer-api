package online_route

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Hello returns a simple health check response.
// @Summary Health check
// @Description Returns "Hello, World!" to confirm the service is running.
// @Tags System
// @Produce text/plain
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func Hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
