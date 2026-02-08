package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite", "file:rbac.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
