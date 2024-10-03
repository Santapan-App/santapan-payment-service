package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	postgresCommands "tobby/internal/repository/postgres/commands"
	postgresQueries "tobby/internal/repository/postgres/queries"
	"tobby/internal/rest"
	pkgEcho "tobby/pkg/echo"
	"tobby/pkg/sql"
	"tobby/token"
	"tobby/user"

	"github.com/joho/godotenv"
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

	userQueryRepo := postgresQueries.NewPostgresUserQueryRepository(conn)
	userQueryCommand := postgresCommands.NewPostgresUserCommandRepository(conn)

	tokenQueryRepo := postgresQueries.NewPostgresTokenQueryRepository(conn)
	tokenCommandRepo := postgresCommands.NewPostgresTokenCommandRepository(conn)

	tokenService := token.NewService(tokenQueryRepo, tokenCommandRepo)
	userService := user.NewService(userQueryRepo, userQueryCommand)

	e := pkgEcho.Setup()

	rest.NewAuthHandler(e, tokenService, userService)

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
