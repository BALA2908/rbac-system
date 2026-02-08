package db

import (
	"database/sql"
	"encoding/json"

	"rbac-backend/internal/models"
)

func GetPermissionsByRole(db *sql.DB, role string) (models.Permissions, error) {
	var raw string
	err := db.QueryRow(
		"SELECT permissions FROM role_permissions WHERE role=?",
		role,
	).Scan(&raw)
	if err != nil {
		return nil, err
	}

	var perms models.Permissions
	err = json.Unmarshal([]byte(raw), &perms)
	return perms, err
}

// ListRoles returns all role names that have config in role_permissions (non-ADMIN roles).
func ListRoles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT role FROM role_permissions ORDER BY role")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// UpdateRolePermissions sets the JSON permissions for a role (MANAGER, EDITOR, VIEWER only).
func UpdateRolePermissions(db *sql.DB, role string, perms models.Permissions) error {
	data, err := json.Marshal(perms)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT OR REPLACE INTO role_permissions (role, permissions) VALUES (?, ?)",
		role, string(data),
	)
	return err
}
