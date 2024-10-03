package sql

import (
	database "database/sql"
	"log"
)

func Setup() *database.DB {
	conn, err := NewConnection()

	if err != nil {
		log.Fatal("Failed to open connection to database", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	return conn
}
