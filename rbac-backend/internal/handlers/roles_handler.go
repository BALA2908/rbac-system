package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"rbac-backend/internal/db"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
)

// RolesHandler handles admin-only role config (get/update permissions from DB).
type RolesHandler struct {
	DB *sql.DB
}

func NewRolesHandler(database *sql.DB) *RolesHandler {
	return &RolesHandler{DB: database}
}

var validConfigRoles = map[string]bool{
	rbac.RoleManager: true,
	rbac.RoleEditor:  true,
	rbac.RoleViewer:  true,
}

func roleFromPath(path string) string {
	const prefix = "/admin/roles/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(path, prefix))
}

func (h *RolesHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roles, err := db.ListRoles(h.DB)
	if err != nil {
		http.Error(w, "failed to list roles", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"roles": roles})
}

func (h *RolesHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	role := roleFromPath(r.URL.Path)
	if role == "" {
		http.Error(w, "role required", http.StatusBadRequest)
		return
	}
	if !validConfigRoles[role] {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}
	perms, err := db.GetPermissionsByRole(h.DB, role)
	if err != nil {
		http.Error(w, "role not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(perms)
}

func (h *RolesHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	role := roleFromPath(r.URL.Path)
	if role == "" {
		http.Error(w, "role required", http.StatusBadRequest)
		return
	}
	if !validConfigRoles[role] {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}
	var perms models.Permissions
	if err := json.NewDecoder(r.Body).Decode(&perms); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := db.UpdateRolePermissions(h.DB, role, perms); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// ServeRoleDetail handles GET (get one) and PUT/POST (update) for /admin/roles/{role}.
func (h *RolesHandler) ServeRoleDetail(w http.ResponseWriter, r *http.Request) {
	role := roleFromPath(r.URL.Path)
	if role == "" {
		http.Error(w, "role required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.GetRole(w, r)
	case http.MethodPut, http.MethodPost:
		h.UpdateRole(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
