package json

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Response(c echo.Context, statusCode int, success bool, message string, data interface{}) error {
	logrus.WithFields(logrus.Fields{
		"statusCode": statusCode,
		"success":    success,
		"message":    message,
	}).Info("Response")

	return c.JSON(statusCode, map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	})
}
