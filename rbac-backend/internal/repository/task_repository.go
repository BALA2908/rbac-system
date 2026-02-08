package repositories

import (
	"database/sql"
	"encoding/json"
	"rbac-backend/internal/models"
	"time"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTask(t models.Task) error {
	// marshal assignees to JSON string stored in assignee column
	var ajson sql.NullString
	if len(t.Assignees) > 0 {
		b, _ := json.Marshal(t.Assignees)
		ajson = sql.NullString{String: string(b), Valid: true}
	}

	_, err := r.DB.Exec(`INSERT INTO tasks (id, project_id, title, description, status, assignee, created_by, started_at, completed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.Title, t.Description, t.Status, ajson, t.CreatedBy, t.StartedAt, t.CompletedAt,
	)
	return err
}

func (r *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	row := r.DB.QueryRow(`SELECT id, project_id, title, description, status, assignee, created_by, started_at, completed_at, created_at, updated_at FROM tasks WHERE id = ?`, id)

	var t models.Task
	var started sql.NullTime
	var completed sql.NullTime

	var astring sql.NullString
	err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &astring, &t.CreatedBy, &started, &completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if started.Valid {
		t.StartedAt = &started.Time
	}
	if completed.Valid {
		t.CompletedAt = &completed.Time
	}

	if astring.Valid && astring.String != "" {
		var arr []string
		if err := json.Unmarshal([]byte(astring.String), &arr); err == nil {
			t.Assignees = arr
		}
	}

	return &t, nil
}

func (r *TaskRepository) ListTasksByProject(projectID string) ([]models.Task, error) {
	rows, err := r.DB.Query(`SELECT id, project_id, title, description, status, assignee, created_by, started_at, completed_at, created_at, updated_at FROM tasks WHERE project_id=? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var started sql.NullTime
		var completed sql.NullTime
		var astring sql.NullString
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &astring, &t.CreatedBy, &started, &completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if started.Valid {
			t.StartedAt = &started.Time
		}
		if completed.Valid {
			t.CompletedAt = &completed.Time
		}
		if astring.Valid && astring.String != "" {
			var arr []string
			if err := json.Unmarshal([]byte(astring.String), &arr); err == nil {
				t.Assignees = arr
			}
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) ListTasksByAssignee(userID string) ([]models.Task, error) {
	// assignee column now stores JSON array; use LIKE to find userID in text
	rows, err := r.DB.Query(`SELECT id, project_id, title, description, status, assignee, created_by, started_at, completed_at, created_at, updated_at FROM tasks WHERE assignee LIKE ? ORDER BY created_at DESC`, "%"+userID+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var started sql.NullTime
		var completed sql.NullTime
		var astring sql.NullString
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &astring, &t.CreatedBy, &started, &completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if started.Valid {
			t.StartedAt = &started.Time
		}
		if completed.Valid {
			t.CompletedAt = &completed.Time
		}
		if astring.Valid && astring.String != "" {
			var arr []string
			if err := json.Unmarshal([]byte(astring.String), &arr); err == nil {
				t.Assignees = arr
			}
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) UpdateTask(t models.Task) error {
	var ajson sql.NullString
	if len(t.Assignees) > 0 {
		b, _ := json.Marshal(t.Assignees)
		ajson = sql.NullString{String: string(b), Valid: true}
	}
	_, err := r.DB.Exec(`UPDATE tasks SET title=?, description=?, status=?, assignee=?, started_at=?, completed_at=?, updated_at=? WHERE id=?`,
		t.Title, t.Description, t.Status, ajson, t.StartedAt, t.CompletedAt, time.Now(), t.ID,
	)
	return err
}

func (r *TaskRepository) AssignTask(taskID, userID string) error {
	// append single userID to assignees array if not already present
	t, err := r.GetTaskByID(taskID)
	if err != nil || t == nil {
		return err
	}
	exists := false
	for _, a := range t.Assignees {
		if a == userID {
			exists = true
			break
		}
	}
	if !exists {
		t.Assignees = append(t.Assignees, userID)
	}
	return r.UpdateTask(*t)
}

func (r *TaskRepository) UpdateStatus(taskID, status string) error {
	// set started_at if moving to IN_PROGRESS and not already set
	if status == "IN_PROGRESS" {
		_, err := r.DB.Exec(`UPDATE tasks SET status=?, started_at = COALESCE(started_at, ?), updated_at=? WHERE id=?`, status, time.Now(), time.Now(), taskID)
		return err
	}
	// set completed_at when marking DONE or ARCHIVED
	if status == "DONE" || status == "ARCHIVED" {
		_, err := r.DB.Exec(`UPDATE tasks SET status=?, completed_at = ?, updated_at=? WHERE id=?`, status, time.Now(), time.Now(), taskID)
		return err
	}
	_, err := r.DB.Exec(`UPDATE tasks SET status=?, updated_at=? WHERE id=?`, status, time.Now(), taskID)
	return err
}

func (r *TaskRepository) DeleteTask(id string) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE id=?`, id)
	return err
}
