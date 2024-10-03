package json

import (
	"github.com/labstack/echo/v4"
)

func Response(c echo.Context, statusCode int, success bool, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	})
}
