package middleware

import (
	"net/http"

	"rbac-backend/internal/rbac"
)

// RequireAdmin restricts the route to ADMIN only (manage users, create/update roles).
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleVal := r.Context().Value(RoleKey)
		if roleVal == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		role := roleVal.(string)
		if role != rbac.RoleAdmin {
			http.Error(w, "admin only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
