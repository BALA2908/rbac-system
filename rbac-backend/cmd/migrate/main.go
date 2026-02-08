package main

import (
	"log"

	"rbac-backend/internal/db"
)

func main() {
	database := db.Connect()
	defer database.Close()

	log.Println("Starting database migrations...")
	if err := db.RunMigrations(database); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("✅ Database migrations completed successfully!")

	rows, _ := database.Query("SELECT role FROM role_permissions")
	for rows.Next() {
		var role string
		rows.Scan(&role)
		log.Println("ROLE:", role)
	}

}
