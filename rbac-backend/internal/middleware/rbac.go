package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"rbac-backend/internal/db"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
)

// contextKey type for RBAC context values.
type contextKey string

const (
	// TablePermKey is the context key for the current table's ResourcePermission (for field-level filtering in handlers).
	TablePermKey contextKey = "tablePerm"
)

// fullAccessPerm returns a ResourcePermission that allows all table and field access (for ADMIN).
func fullAccessPerm() models.ResourcePermission {
	return models.ResourcePermission{
		View:   true,
		Create: true,
		Edit:   true,
		Delete: true,
		Fields: nil, // nil => allow all fields in utils.FilterFields / FilterEditableFields
	}
}

// RBACMiddleware enforces config-driven RBAC: ADMIN has full access; other roles use DB config only.
func RBACMiddleware(database *sql.DB, table, action string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleVal := r.Context().Value(RoleKey)
		if roleVal == nil {
			http.Error(w, "role missing in context", http.StatusUnauthorized)
			return
		}
		role := roleVal.(string)

		// Enforce users table is ADMIN-only regardless of DB config.
		if table == "users" && role != rbac.RoleAdmin {
			http.Error(w, "users table restricted to ADMIN", http.StatusForbidden)
			return
		}

		var tablePerm models.ResourcePermission

		if role == rbac.RoleAdmin {
			// ADMIN: full access to all tables and fields; no DB lookup.
			tablePerm = fullAccessPerm()
		} else {
			// Other roles: permissions ONLY from DB (no hardcoded logic).
			perms, err := db.GetPermissionsByRole(database, role)
			if err != nil {
				http.Error(w, "permission lookup failed", http.StatusForbidden)
				return
			}

			var ok bool
			tablePerm, ok = perms[table]
			if !ok {
				http.Error(w, "no table access", http.StatusForbidden)
				return
			}

			switch action {
			case rbac.ActionView:
				if !tablePerm.View {
					http.Error(w, "view not allowed", http.StatusForbidden)
					return
				}
			case rbac.ActionCreate:
				if !tablePerm.Create {
					http.Error(w, "create not allowed", http.StatusForbidden)
					return
				}
			case rbac.ActionEdit:
				if !tablePerm.Edit {
					http.Error(w, "edit not allowed", http.StatusForbidden)
					return
				}
			case rbac.ActionDelete:
				if !tablePerm.Delete {
					http.Error(w, "delete not allowed", http.StatusForbidden)
					return
				}
			default:
				http.Error(w, "unknown action", http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), TablePermKey, tablePerm)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
