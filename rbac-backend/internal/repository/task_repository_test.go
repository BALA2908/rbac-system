package repositories

import (
	"database/sql"
	"testing"

	"rbac-backend/internal/models"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?_foreign_keys=on")
	if err != nil {
		t.Fatal(err)
	}

	// create minimal tables needed for tasks
	schema := `
    CREATE TABLE users (id TEXT PRIMARY KEY, name TEXT);
    CREATE TABLE projects (id TEXT PRIMARY KEY, name TEXT, created_by TEXT);
    CREATE TABLE tasks (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, title TEXT NOT NULL, description TEXT, status TEXT NOT NULL DEFAULT 'TODO', assignee TEXT, created_by TEXT NOT NULL, started_at DATETIME, completed_at DATETIME, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);
    `
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCreateAndGetTask(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := models.Task{ID: "tid1", ProjectID: "pid1", Title: "Test", CreatedBy: "u1", Status: "TODO"}

	if err := repo.CreateTask(task); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := repo.GetTaskByID("tid1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got == nil || got.ID != "tid1" || got.Title != "Test" {
		t.Fatalf("unexpected task: %+v", got)
	}
}
