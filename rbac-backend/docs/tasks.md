# Tasks API

Endpoints (requires Authorization: `Bearer <token>`):

- `POST /tasks/create` — create a task. JSON body: `{ "project_id": "<id>", "title": "...", "description": "...", "assignees": ["<user_id>"] }`.
- `GET /tasks?project_id=<id>` — list tasks for a project.
- `GET /tasks?assignee=<user_id>` — list tasks assigned to a user (matches any id present in `assignees`).
- `GET /tasks/get?id=<id>` — get single task.
- `POST /tasks/update` — update task. JSON body must include `id` and editable fields (title, description, status, assignees).
- `POST /tasks/assign` — assign task. JSON body: `{ "id": "<task_id>", "assignees": ["<user_id>"] }` or `{ "id": "<task_id>", "assignee": "<user_id>" }` to append a single assignee.
- `GET /tasks/delete?id=<id>` — delete task (protected by RBAC delete permission).

Status flow: `TODO -> IN_PROGRESS -> REVIEW -> DONE`. Handlers set timestamps when starting or completing.

Frontend:
- `rbac-frontend/src/pages/Tasks.tsx` — simple page with create form.
- `rbac-frontend/src/components/KanbanBoard.tsx` — Kanban UI.
- `rbac-frontend/src/api/tasks.ts` — API helpers (reads token from `localStorage.token`).
