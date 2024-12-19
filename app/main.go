package main

import (
	"log"
	"os"
	"os/signal"
	"santapan_payment_service/internal/rest"
	"santapan_payment_service/payment"
	pkgEcho "santapan_payment_service/pkg/echo"

	postgresCommands "santapan_payment_service/internal/repository/postgres/commands"
	postgresQueries "santapan_payment_service/internal/repository/postgres/queries"

	"santapan_payment_service/pkg/sql"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the postgres driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import the file source driver

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	conn := sql.Setup()
	defer sql.Close(conn)

	// Initialize services

	e := pkgEcho.Setup()

	paymentQuery := postgresQueries.NewPostgresPaymentQueryRepository(conn)
	paymentCommand := postgresCommands.NewPostgresPaymentCommandRepository(conn)

	paymentService := payment.NewService(paymentQuery, paymentCommand)
	// Initialize the REST API
	rest.NewPaymentHandler(e, paymentService)
	go func() {
		pkgEcho.Start(e)
	}()

	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-quit

	pkgEcho.Shutdown(e, defaultTimeout)
}
