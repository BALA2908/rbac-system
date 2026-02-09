package repositories

import (
	"database/sql"
	"errors"
	"rbac-backend/internal/models"
	"strings"
)

type ProjectRepository struct {
	DB *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}
func (r *ProjectRepository) CreateProject(project models.Project) error {

	_, err := r.DB.Exec(
		`INSERT INTO projects (id, name, description, created_by)
		 VALUES (?, ?, ?, ?)`,
		project.ID,
		project.Name,
		project.Description,
		project.CreatedBy,
	)

	return err
}
func (r *ProjectRepository) CreateProjectDynamic(data map[string]interface{}) error {

	if len(data) == 0 {
		return errors.New("no data provided")
	}

	columns := []string{}
	placeholders := []string{}
	args := []interface{}{}

	for col, val := range data {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		args = append(args, val)
	}

	_ = "INSERT INTO projects (" +
		strings.Join(columns, ",") +
		") VALUES (" +
		strings.Join(placeholders, ",") +
		")"

	var assignments []string
	if a, ok := data["assigned_employees"]; ok {
		switch v := a.(type) {
		case []string:
			assignments = v
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					assignments = append(assignments, s)
				}
			}
		}

		newCols := []string{}
		newPlaceholders := []string{}
		newArgs := []interface{}{}
		for i, col := range columns {
			if col == "assigned_employees" {
				continue
			}
			newCols = append(newCols, col)
			newPlaceholders = append(newPlaceholders, placeholders[i])
			newArgs = append(newArgs, args[i])
		}
		columns = newCols
		placeholders = newPlaceholders
		args = newArgs
	}

	_, err := r.DB.Exec("INSERT INTO projects ("+strings.Join(columns, ",")+") VALUES ("+strings.Join(placeholders, ",")+")", args...)
	if err != nil {
		return err
	}

	if len(assignments) > 0 {
		pid, _ := data["id"].(string)
		if pid == "" {
			return errors.New("project id required for assignments")
		}
		tx, err := r.DB.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT OR REPLACE INTO project_assignments (project_id, user_id) VALUES (?, ?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		for _, uid := range assignments {
			if _, err := stmt.Exec(pid, uid); err != nil {
				tx.Rollback()
				return err
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectRepository) GetProjects() ([]models.Project, error) {

	rows, err := r.DB.Query(`SELECT id, name, description, created_by FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy)
		if err != nil {
			return nil, err
		}

		assignRows, err := r.DB.Query(`SELECT user_id FROM project_assignments WHERE project_id = ?`, p.ID)
		if err == nil {
			defer assignRows.Close()
			var assigns []string
			for assignRows.Next() {
				var uid string
				if err := assignRows.Scan(&uid); err == nil {
					assigns = append(assigns, uid)
				}
			}
			p.AssignedEmployees = assigns
		}

		projects = append(projects, p)
	}

	return projects, nil
}
func (r *ProjectRepository) UpdateProjectDynamic(data map[string]interface{}) error {

	idVal, ok := data["id"]
	if !ok {
		return errors.New("id required for update")
	}
	id := idVal.(string)

	delete(data, "id") // do not update ID

	if len(data) == 0 {
		return errors.New("no editable fields provided")
	}

	query := "UPDATE projects SET "
	args := []interface{}{}
	i := 0

	for field, value := range data {
		if i > 0 {
			query += ", "
		}
		query += field + "=?"
		args = append(args, value)
		i++
	}

	query += " WHERE id=?"
	args = append(args, id)

	_, err := r.DB.Exec(query, args...)
	return err
}

func (r *ProjectRepository) DeleteProject(id string) error {

	_, err := r.DB.Exec(`DELETE FROM projects WHERE id=?`, id)
	return err
}
