package db

import (
	"database/sql"
	"fmt"
	"log"

	"rbac-backend/internal/auth"

	"github.com/google/uuid"
)

func SeedAdmin(db *sql.DB) {
	// Admin details
	adminID := uuid.New().String()
	name := "Super Admin"
	email := "admin@example.com"
	password := "admin123" // you will use this to login
	role := "ADMIN"

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Insert admin (ignore if already exists)
	query := `
	INSERT OR IGNORE INTO users 
	(id, name, email, password_hash, role, is_active)
	VALUES (?, ?, ?, ?, ?, 1)
	`

	_, err = db.Exec(query, adminID, name, email, hashedPassword, role)
	if err != nil {
		log.Fatal("Failed to seed admin:", err)
	}

	fmt.Println("✅ Admin user seeded successfully")
	fmt.Println("📧 Email:", email)
	fmt.Println("🔑 Password:", password)
}
