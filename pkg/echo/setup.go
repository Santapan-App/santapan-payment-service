package echo

import (
	"log"
	"os"
	"strconv"
	"time"

	localMiddleware "santapan_payment_service/pkg/echo/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9091"
)

func Setup() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(localMiddleware.CORS)

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("Failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}

	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(localMiddleware.SetRequestContextWithTimeout(timeoutContext))

	return e
}
