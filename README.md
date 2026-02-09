# RBAC System - Config-Driven Access Control

A full-stack Role-Based Access Control (RBAC) system that demonstrates table-level and field-level permissions managed dynamically by an Admin user.

## Project Overview

This is a comprehensive RBAC system built with:
- **Backend:** Golang + SQLite
- **Frontend:** React with TypeScript + Vite
- **Database:** SQLite with migrations

The system enforces permissions at two levels:
1. **Table-level**: Control which tables users can view/edit
2. **Field-level**: Control which columns users can view/edit within those tables

All permissions are **configuration-driven** (stored in database, not hardcoded) and managed by administrators.

## Role & Permission Model

### Built-in Roles
- **Admin**: Full system access, can manage roles and permissions
- **Manager**: Can view & edit projects and tasks
- **Viewer**: Can only view data (read-only access)
- **Editor**: Can edit projects but only view tasks

### Entities (Sample Data)
1. **Users** - System users with assigned roles
2. **Projects** - Organizational projects
3. **Tasks** - Project tasks with status tracking (TODO → IN_PROGRESS → REVIEW → DONE)

### Permission Enforcement
- ✅ Unauthorized tables are hidden from the UI
- ✅ Restricted fields are not visible in the frontend
- ✅ Edit/Delete actions require backend permission validation
- ✅ Unauthorized requests are rejected with 403 Forbidden

## Quick Start

### Prerequisites
- Go 1.25+
- Node.js 16+
- SQLite (included with Go)

### Backend Setup

```bash
cd rbac-backend

# Copy environment configuration
cp .env.example .env

# Update .env with your configuration (change JWT_SECRET in production)
# JWT_SECRET=your-secret-key-here
# PORT=8080
# DB_PATH=rbac.db

# Initialize database & run migrations
go run cmd/migrate/main.go

# (Optional) Seed sample data
go run cmd/seed/main.go

# Start the server
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

### Frontend Setup

```bash
cd rbac-frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend runs on `http://localhost:5173`

## Configuration

Create a `.env` file in the `rbac-backend` directory (use `.env.example` as template):

```env
JWT_SECRET=your-secret-key-here-change-in-production
PORT=8080
DB_PATH=rbac.db
```


### Admin Account for Login
- **Email:** "admin@example.com"
- **Password:** admin123
  **Role:** Admin

### Regular User Accounts
- **Username:** managerf@rbac.com
  **Password:** 123456
  **Role:** Manager

- **Email:** "editor@example.com"
- **Password:** 123456
  **Role:** Editor

- **Email:** "viewer@example.com"
- **Password:** 123456
  **Role:** Viewer


## Architecture

```
rbac-backend/
├── internal/
│   ├── auth/         # JWT & password utilities
│   ├── db/           # Database & migrations
│   ├── handlers/     # API endpoint handlers
│   ├── middleware/   # Auth & RBAC checks
│   ├── models/       # Data structures
│   ├── rbac/         # Permission logic
│   └── repository/   # Database queries
├── migrations/       # SQL migration files
└── cmd/
    ├── server/       # Main server entry point
    └── migrate/      # Database migration runner

rbac-frontend/
├── src/
│   ├── components/   # React components
│   ├── pages/        # Page components
│   ├── api/          # API client code
│   └── utils/        # Helper functions
```

Try performing restricted actions (edit as Viewer) to see 403 errors from backend.

