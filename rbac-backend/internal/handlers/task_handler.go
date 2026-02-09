package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
	repositories "rbac-backend/internal/repository"
	"rbac-backend/internal/utils"
)

type TaskHandler struct {
	Repo *repositories.TaskRepository
}

func NewTaskHandler(repo *repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	pid, ok := incoming["project_id"].(string)
	if !ok || pid == "" {
		http.Error(w, "project_id required", http.StatusBadRequest)
		return
	}
	title, ok := incoming["title"].(string)
	if !ok || title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	t := models.Task{
		ID:        uuid.New().String(),
		ProjectID: pid,
		Title:     title,
		CreatedBy: userID,
		Status:    "TODO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if desc, ok := safe["description"].(string); ok {
		t.Description = desc
	}
	if s, ok := safe["status"].(string); ok && s != "" {
		t.Status = s
	}
	if arr, ok := safe["assignees"].([]interface{}); ok {
		for _, it := range arr {
			if sid, ok := it.(string); ok {
				t.Assignees = append(t.Assignees, sid)
			}
		}
	} else if a, ok := safe["assignee"].(string); ok && a != "" {
		t.Assignees = append(t.Assignees, a)
	}

	if err := h.Repo.CreateTask(t); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	projectID := r.URL.Query().Get("project_id")
	assignee := r.URL.Query().Get("assignee")

	var tasks []models.Task
	var err error

	if projectID != "" {
		tasks, err = h.Repo.ListTasksByProject(projectID)
	} else if assignee != "" {
		tasks, err = h.Repo.ListTasksByAssignee(assignee)
	} else {
		http.Error(w, "project_id or assignee query required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	var out []map[string]interface{}
	for _, t := range tasks {
		if role == rbac.RoleViewer {
			allowed := t.CreatedBy == userID
			if !allowed {
				for _, a := range t.Assignees {
					if a == userID {
						allowed = true
						break
					}
				}
			}
			if !allowed {
				continue
			}
		}

		row := map[string]interface{}{
			"id":           t.ID,
			"project_id":   t.ProjectID,
			"title":        t.Title,
			"description":  t.Description,
			"status":       t.Status,
			"assignees":    t.Assignees,
			"created_by":   t.CreatedBy,
			"started_at":   t.StartedAt,
			"completed_at": t.CompletedAt,
			"created_at":   t.CreatedAt,
			"updated_at":   t.UpdatedAt,
		}

		filtered := utils.FilterFields(row, tablePerm.Fields)
		out = append(out, filtered)
	}

	json.NewEncoder(w).Encode(out)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	t, err := h.Repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "failed to fetch task", http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	if role != rbac.RoleAdmin && role != rbac.RoleManager && role != rbac.RoleEditor {
		allowed := t.CreatedBy == userID
		if !allowed {
			for _, a := range t.Assignees {
				if a == userID {
					allowed = true
					break
				}
			}
		}
		if !allowed {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	row := map[string]interface{}{
		"id":           t.ID,
		"project_id":   t.ProjectID,
		"title":        t.Title,
		"description":  t.Description,
		"status":       t.Status,
		"assignees":    t.Assignees,
		"created_by":   t.CreatedBy,
		"started_at":   t.StartedAt,
		"completed_at": t.CompletedAt,
		"created_at":   t.CreatedAt,
		"updated_at":   t.UpdatedAt,
	}

	json.NewEncoder(w).Encode(utils.FilterFields(row, tablePerm.Fields))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	idVal, ok := incoming["id"].(string)
	if !ok || idVal == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	existing, err := h.Repo.GetTaskByID(idVal)
	if err != nil || existing == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	if title, ok := safe["title"].(string); ok {
		existing.Title = title
	}
	if desc, ok := safe["description"].(string); ok {
		existing.Description = desc
	}
	if status, ok := safe["status"].(string); ok {
		existing.Status = status
		if status == "IN_PROGRESS" && existing.StartedAt == nil {
			now := time.Now()
			existing.StartedAt = &now
		}
		if (status == "DONE" || status == "ARCHIVED") && existing.CompletedAt == nil {
			now := time.Now()
			existing.CompletedAt = &now
		}
	}
	if arr, ok := safe["assignees"].([]interface{}); ok {
		existing.Assignees = nil
		for _, it := range arr {
			if sid, ok := it.(string); ok {
				existing.Assignees = append(existing.Assignees, sid)
			}
		}
	} else if a, ok := safe["assignee"].(string); ok && a != "" {
		existing.Assignees = []string{a}
	}

	existing.UpdatedAt = time.Now()

	if err := h.Repo.UpdateTask(*existing); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		ID        string   `json:"id"`
		Assignee  string   `json:"assignee"`
		Assignees []string `json:"assignees"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if payload.ID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if len(payload.Assignees) > 0 {
		t, err := h.Repo.GetTaskByID(payload.ID)
		if err != nil || t == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		t.Assignees = payload.Assignees
		if err := h.Repo.UpdateTask(*t); err != nil {
			http.Error(w, "assign failed", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "assigned"})
		return
	}

	if payload.Assignee == "" {
		http.Error(w, "assignee required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.AssignTask(payload.ID, payload.Assignee); err != nil {
		http.Error(w, "assign failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteTask(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "task deleted"})
}
