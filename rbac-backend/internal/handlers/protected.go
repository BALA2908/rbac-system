package handlers

import (
	"net/http"

	"rbac-backend/internal/middleware"
)

func Protected() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey).(string)
		role := r.Context().Value(middleware.RoleKey).(string)

		w.Write([]byte("Access granted ✅ UserID: " + userID + " Role: " + role))
	}
}
