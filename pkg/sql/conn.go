package sql

import (
	database "database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

func NewConnection() (*database.DB, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	// Construct the connection string with port
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("sslmode", "disable") // Disable SSL
	val.Add("timezone", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	// Open database connection
	dbConn, err := database.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
