// internal/models/permission.go
package models

// FieldPermission defines field-level access (view / create / edit).
// Used in JSON config in role_permissions.permissions.
type FieldPermission struct {
	View   bool `json:"view"`
	Create bool `json:"create"`
	Edit   bool `json:"edit"`
}

// ResourcePermission defines table-level and optional field-level permissions.
// Permissions for other roles come ONLY from DB; ADMIN is handled in code.
type ResourcePermission struct {
	View   bool                       `json:"view"`
	Create bool                       `json:"create"`
	Edit   bool                       `json:"edit"`
	Delete bool                       `json:"delete"`
	Fields map[string]FieldPermission `json:"fields,omitempty"`
}

// Permissions is keyed by table/resource name (e.g. "projects", "users").
type Permissions map[string]ResourcePermission
