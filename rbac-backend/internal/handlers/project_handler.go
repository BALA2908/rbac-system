package handlers

import (
	"encoding/json"
	"net/http"
	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
	repositories "rbac-backend/internal/repository"
	"rbac-backend/internal/utils"

	"github.com/google/uuid"
)

type ProjectHandler struct {
	Repo *repositories.ProjectRepository
}

func NewProjectHandler(repo *repositories.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{Repo: repo}
}
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	safe["id"] = uuid.New().String()
	safe["created_by"] = userID

	err := h.Repo.CreateProjectDynamic(safe)
	if err != nil {
		http.Error(w, "failed to create project", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(safe)
}

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	projects, err := h.Repo.GetProjects()
	if err != nil {
		http.Error(w, "failed to fetch projects", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}

	// Determine role and userID to apply project-level visibility rules
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	for _, p := range projects {
		// If not ADMIN, only include projects assigned to the current user
		if role != rbac.RoleAdmin {
			allowed := false
			for _, uid := range p.AssignedEmployees {
				if uid == userID {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}
		row := map[string]interface{}{
			"id":                 p.ID,
			"name":               p.Name,
			"description":        p.Description,
			"created_by":         p.CreatedBy,
			"assigned_employees": p.AssignedEmployees,
		}

		filtered := utils.FilterFields(row, tablePerm.Fields)

		// Ensure assigned_employees is included when the table view is allowed.
		// Field-level config may omit this field; include it here so frontend can
		// show assignees when the user can view projects at all.
		if len(p.AssignedEmployees) > 0 && tablePerm.View {
			filtered["assigned_employees"] = p.AssignedEmployees
		}
		response = append(response, filtered)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// ✅ ID must be read before filtering
	id, ok := incoming["id"].(string)
	if !ok || id == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	// Remove ID from editable fields input
	delete(incoming, "id")

	// Filter only editable fields
	safeData := utils.FilterEditableFields(incoming, tablePerm.Fields)

	if len(safeData) == 0 {
		http.Error(w, "no editable fields", http.StatusForbidden)
		return
	}

	// Put ID back for repository WHERE clause
	safeData["id"] = id

	err := h.Repo.UpdateProjectDynamic(safeData)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	return // ✅ VERY IMPORTANT

}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	err := h.Repo.DeleteProject(id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "project deleted",
	})
}
