package rbac

// Role names. ADMIN is special (full access in code); others use DB config only.
const (
	RoleAdmin   = "ADMIN"
	RoleManager = "MANAGER"
	RoleEditor  = "EDITOR"
	RoleViewer  = "VIEWER"
)

// Table/resource names used in permission config.
const (
	TableProjects = "projects"
	TableUsers    = "users"
)

// Actions for table-level checks.
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionEdit   = "edit"
	ActionDelete = "delete"
)
