package sql

import (
	database "database/sql"
	"log"
)

func Close(conn *database.DB) {
	if err := conn.Close(); err != nil {
		log.Fatal("Error when closing the DB connection", err)
	}
}
