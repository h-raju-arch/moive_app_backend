package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Open() *sql.DB {
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		log.Fatal("Empty database string")
	}
	var err error
	db, err := sql.Open("postgres", db_url)

	if err != nil {
		log.Fatal("Error while Opening db connection", err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal("Error while pinging db", pingErr)
	}
	return db
}
