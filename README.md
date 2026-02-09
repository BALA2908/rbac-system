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

**Important:** 
- ✅ Add `.env` to `.gitignore` (contains secrets)
- ⚠️ Change `JWT_SECRET` to a strong random key in production
- `.env.example` is included for reference

### Admin Account
- **Username:** `admin`
- **Password:** `password123`

### Regular User Accounts
- **Username:** `john_doe`  
  **Password:** `password123`  
  **Role:** Manager

- **Username:** `jane_smith`  
  **Password:** `password123`  
  **Role:** Editor

- **Username:** `viewer_user`  
  **Password:** `password123`  
  **Role:** Viewer

## Key Features

### Admin Dashboard
- Create and manage roles
- Define permissions per role
- View all users and their assigned roles
- Real-time permission enforcement

### User Dashboard
- View only accessible tables
- See only accessible fields
- Perform actions based on permissions
- Kanban board for task management

### API Endpoints

**Authentication**
- `POST /auth/login` - Login with username/password
- `POST /auth/register` - Register new user (admin only)

**Admin Features**
- `GET /admin/users` - List all users
- `POST /admin/users` - Create user
- `GET /admin/roles` - List all roles
- `POST /admin/roles` - Create role
- `POST /admin/permissions` - Grant permission to role

**User Data**
- `GET /projects` - List projects (respects permissions)
- `POST /projects` - Create project (if permitted)
- `GET /tasks` - List tasks with permission filtering
- `POST /tasks/create` - Create task

All endpoints require JWT token in `Authorization: Bearer <token>` header.

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

## How It Works

1. **User logs in** with username & password
2. **Backend validates** credentials and issues JWT token
3. **Frontend stores** token in localStorage
4. **User accesses dashboard** with only permitted data
5. **Backend enforces** all permission checks via RBAC middleware
6. **Frontend hides** restricted tables & fields based on role

## Testing Permission Enforcement

1. Login as Admin → can access all features and manage permissions
2. Login as Manager → can see Projects & Tasks, edit both
3. Login as Editor → can see Projects & edit, only view Tasks
4. Login as Viewer → can only view all data (read-only)

Try performing restricted actions (edit as Viewer) to see 403 errors from backend.

## Development Notes

- All permissions are checked at the API level (backend validation is authoritative)
- Frontend UI respects permissions for better UX but doesn't rely on it for security
- Database migrations are run automatically on server startup
- JWT tokens expire after configured duration
- SQLite database file: `rbac.db`

## License

MIT
