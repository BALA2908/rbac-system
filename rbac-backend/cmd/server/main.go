package main

import (
	"log"
	"net/http"

	"rbac-backend/internal/config"
	"rbac-backend/internal/db"
	"rbac-backend/internal/handlers"
	"rbac-backend/internal/middleware"
	repositories "rbac-backend/internal/repository"
)

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load configuration from .env file
	config.LoadConfig()

	database := db.Connect()
	defer database.Close()

	// Root Greeting
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to RBAC System Backend 🚀"))
	})

	// AUTH ROUTES
	// POST /login - Authenticate user with username/password and return JWT token
	http.Handle("/login", handlers.Login(database))

	// ⭐ CREATE PROJECT HANDLER
	projectRepo := repositories.NewProjectRepository(database)
	projectHandler := handlers.NewProjectHandler(projectRepo)
	userRepo := repositories.NewUserRepository(database)
	adminHandler := handlers.NewAdminHandler(userRepo)

	// TASK HANDLER
	taskRepo := repositories.NewTaskRepository(database)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// ⭐ PROJECT ROUTES
	// POST /projects/create - Create a new project (requires create permission)
	http.Handle(
		"/projects/create",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "create",
				http.HandlerFunc(projectHandler.CreateProject),
			),
		),
	)

	// GET /projects - List all projects user has access to (requires view permission)
	http.Handle(
		"/projects",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "view",
				http.HandlerFunc(projectHandler.GetProjects),
			),
		),
	)

	// PUT /projects/update - Update an existing project (requires edit permission)
	http.Handle(
		"/projects/update",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "edit",
				http.HandlerFunc(projectHandler.UpdateProject),
			),
		),
	)

	// DELETE /projects/delete - Delete a project (requires create permission)
	http.Handle(
		"/projects/delete",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "create",
				http.HandlerFunc(projectHandler.DeleteProject),
			),
		),
	)
	// POST /tasks/create - Create a new task (requires create permission)
	http.Handle(
		"/tasks/create",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "create",
				http.HandlerFunc(taskHandler.CreateTask),
			),
		),
	)

	// GET /tasks - List all tasks with permission filtering (requires view permission)
	http.Handle(
		"/tasks",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "view",
				http.HandlerFunc(taskHandler.ListTasks),
			),
		),
	)

	// GET /tasks/get - Get a single task by ID (requires view permission)
	http.Handle(
		"/tasks/get",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "view",
				http.HandlerFunc(taskHandler.GetTask),
			),
		),
	)

	// PUT /tasks/update - Update task details like title, description, status (requires edit permission)
	http.Handle(
		"/tasks/update",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "edit",
				http.HandlerFunc(taskHandler.UpdateTask),
			),
		),
	)

	// POST /tasks/assign - Assign task to users (requires edit permission)
	http.Handle(
		"/tasks/assign",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "edit",
				http.HandlerFunc(taskHandler.AssignTask),
			),
		),
	)

	// DELETE /tasks/delete - Delete a task (requires delete permission)
	http.Handle(
		"/tasks/delete",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "delete",
				http.HandlerFunc(taskHandler.DeleteTask),
			),
		),
	)

	// ADMIN ROUTES
	// POST /admin/create-user - Create a new user with assigned role (admin only, requires edit permission)
	http.Handle(
		"/admin/create-user",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "edit",
				http.HandlerFunc(adminHandler.CreateUser),
			),
		),
	)

	// GET /api/users - List all users with their roles and permissions (admin only, requires view permission)

	// LIST USERS - protected route
	http.Handle(
		"/api/users",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "view",
				http.HandlerFunc(adminHandler.ListUsers),
			),
		),
	)
	db.SeedAdmin(database)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", CORSMiddleware(http.DefaultServeMux)))
}
