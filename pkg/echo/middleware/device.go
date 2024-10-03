package echo

import (
	"log"
	"tobby/domain"

	"github.com/labstack/echo/v4"
)

func DeviceInfoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract device information from headers
			deviceInfo := domain.DeviceHeaderInformation{
				DeviceID:    c.Request().Header.Get("deviceid"),
				DeviceName:  c.Request().Header.Get("devicename"),
				DeviceBrand: c.Request().Header.Get("devicebrand"),
				DeviceModel: c.Request().Header.Get("devicemodel"),
				IPAddress:   c.Request().Header.Get("ipaddress"),
			}

			log.Default().Println("Device Info: ", deviceInfo)
			// Store device info in context
			c.Set("deviceInfo", deviceInfo)

			// Continue with the next handler
			return next(c)
		}
	}
}
