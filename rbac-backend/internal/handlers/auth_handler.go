package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rbac-backend/internal/auth"
	dbrepo "rbac-backend/internal/db"
	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
	"rbac-backend/internal/utils"

	"github.com/google/uuid"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Email    string
			Password string
		}

		json.NewDecoder(r.Body).Decode(&req)

		var userID, hash, role string
		err := db.QueryRow(
			"SELECT id, password_hash, role FROM users WHERE email=? AND is_active=1",
			req.Email,
		).Scan(&userID, &hash, &role)

		if err != nil || auth.CheckPassword(hash, req.Password) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateJWT(userID, role)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// Signup registers a new user and returns a JWT token.
func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "Name, email and password are required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}

		userID := uuid.New().String()
		role := "VIEWER"

		_, err = db.Exec(
			`INSERT INTO users (id, name, email, password_hash, role, is_active) 
			 VALUES (?, ?, ?, ?, ?, 1)`,
			userID, req.Name, req.Email, hashedPassword, role,
		)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusBadRequest)
			return
		}

		token, err := auth.GenerateJWT(userID, role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func ViewEmployees(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value(middleware.RoleKey).(string)

		perms, _ := dbrepo.GetPermissionsByRole(database, role)
		fieldPerms := perms["employees"].Fields

		employee := map[string]interface{}{
			"id":     "E101",
			"name":   "Ravi",
			"salary": 90000,
		}

		response := utils.FilterFields(employee, fieldPerms)
		json.NewEncoder(w).Encode(response)
	})
}

func EditEmployees(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value(middleware.RoleKey).(string)

		perms, _ := dbrepo.GetPermissionsByRole(database, role)
		fieldPerms := perms["employees"].Fields

		var input map[string]interface{}
		json.NewDecoder(r.Body).Decode(&input)

		allowedUpdate := utils.FilterEditableFields(input, fieldPerms)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"updated_fields": allowedUpdate,
		})
	})
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email, password required", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = rbac.RoleViewer
	}
	validRole := map[string]bool{
		rbac.RoleAdmin: true, rbac.RoleManager: true,
		rbac.RoleEditor: true, rbac.RoleViewer: true,
	}
	if !validRole[req.Role] {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "password error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashed,
		Role:         req.Role,
		IsActive:     true,
	}
	if err := h.UserRepo.CreateUser(user); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "user created"})
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.UserRepo.ListUsers()
	if err != nil {
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}
