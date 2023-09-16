package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDb() error {
	dbFilePath := os.Getenv("DB_FILE")
	// Dev
	if dbFilePath == "" {
		dbFilePath = "./foo.db"
	}
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	Db = db
	return nil
}
