package rbac

import (
	"database/sql"
	"encoding/json"
)

type Permissions map[string]any

func LoadPermissions(db *sql.DB, role string) (Permissions, error) {
	var raw string

	err := db.QueryRow(
		"SELECT permissions FROM role_permissions WHERE role = ?",
		role,
	).Scan(&raw)
	if err != nil {
		return nil, err
	}

	var perms Permissions
	err = json.Unmarshal([]byte(raw), &perms)
	return perms, err
}
