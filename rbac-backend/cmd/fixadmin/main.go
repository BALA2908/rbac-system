package main

import (
	"log"
	"rbac-backend/internal/db"
)

func main() {
	database := db.Connect()
	defer database.Close()

	log.Println("Updating ADMIN permissions...")

	_, err := database.Exec(`
DELETE FROM role_permissions WHERE role='ADMIN';

INSERT INTO role_permissions (role, permissions) VALUES (
'ADMIN',
'{
  "projects": {
    "view": true,
    "edit": true,
    "fields": {
      "id": { "view": true, "edit": true },
      "name": { "view": true, "edit": true },
      "description": { "view": true, "edit": true },
      "created_by": { "view": true, "edit": true }
    }
  },
  "users": {
    "view": true,
    "edit": true,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": true },
      "email": { "view": true, "edit": true },
      "role": { "view": true, "edit": true },
      "is_active": { "view": true, "edit": true }
    }
  }
}'
);
`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ ADMIN permissions fixed!")
}
