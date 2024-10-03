package echo

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func Start(e *echo.Echo) {
	addrs := os.Getenv("SERVER_ADDRESS")
	log.Default().Println("Starting server at", addrs)
	if addrs == "" {
		addrs = defaultAddress
	}

	if err := e.Start(addrs); err != nil {
		log.Fatal("Shutting down the server")
	}
}

// Shutdown attempts to gracefully shut down the Echo server.
func Shutdown(e *echo.Echo, timeout int) {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	// Create a context with a deadline to wait for ongoing requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
