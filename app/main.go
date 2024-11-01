package main

import (
	"log"
	"os"
	"os/signal"
	"santapan/article"
	"santapan/category"
	postgresCommands "santapan/internal/repository/postgres/commands"
	postgresQueries "santapan/internal/repository/postgres/queries"
	"santapan/internal/rest"
	pkgEcho "santapan/pkg/echo"
	"santapan/pkg/sql"
	"santapan/token"
	"santapan/user"
	"syscall"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
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

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup repositories and services
	userQueryRepo := postgresQueries.NewPostgresUserQueryRepository(conn)
	userQueryCommand := postgresCommands.NewPostgresUserCommandRepository(conn)

	tokenQueryRepo := postgresQueries.NewPostgresTokenQueryRepository(conn)
	tokenCommandRepo := postgresCommands.NewPostgresTokenCommandRepository(conn)

	articleQueryRepo := postgresQueries.NewArticleRepository(conn)
	articleCommandRepo := postgresQueries.NewArticleRepository(conn)

	categoryQueryRepo := postgresQueries.NewCategoryRepository(conn)
	categoryCommandRepo := postgresQueries.NewCategoryRepository(conn)

	tokenService := token.NewService(tokenQueryRepo, tokenCommandRepo)
	userService := user.NewService(userQueryRepo, userQueryCommand)
	articleService := article.NewService(articleQueryRepo, articleCommandRepo)
	categoryService := category.NewService(categoryQueryRepo, categoryCommandRepo)

	e := pkgEcho.Setup()

	rest.NewAuthHandler(e, tokenService, userService)
	rest.NewArticleHandler(e, articleService)
	rest.NewCategoryHandler(e, categoryService)
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

// runMigrations runs the database migrations
func runMigrations() error {
	// Build the database connection string from environment variables
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")

	// Format the connection string
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	// Create a new migration instance
	m, err := migrate.New(
		"app/migrations", // Update the path to your migrations directory
		connectionString,
	)
	if err != nil {
		return err
	}

	// Perform the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
