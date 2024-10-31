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
