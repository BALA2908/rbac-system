package main

import (
	"fmt"
	"log"

	"rbac-backend/internal/db"
)

func main() {
	database := db.Connect()
	defer database.Close()

	// List all tables
	rows, err := database.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%' 
		ORDER BY name
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Tables in rbac.db:")
	fmt.Println("-------------------")

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(tables) == 0 {
		fmt.Println("(no user tables found)")
		return
	}

	for _, t := range tables {
		var count int
		_ = database.QueryRow("SELECT COUNT(*) FROM \"" + t + "\"").Scan(&count)
		fmt.Printf("  %-20s  (%d rows)\n", t, count)
	}
}
