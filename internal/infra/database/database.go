package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDatabaseConn(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic("Error opening database connection", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic("Error pinging database", err)
	}

	log.Println("Database connection established")
	return db
}
