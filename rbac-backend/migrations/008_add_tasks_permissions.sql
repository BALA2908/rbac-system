DELETE FROM role_permissions;

INSERT INTO role_permissions (role, permissions) VALUES
('MANAGER', '{
  "projects": {
    "view": true,
    "create": true,
    "edit": true,
    "delete": true,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": true },
      "description": { "view": true, "edit": true },
      "created_by": { "view": true, "edit": false }
    }
  },
  "tasks": {
    "view": true,
    "create": true,
    "edit": true,
    "delete": true,
    "fields": {
      "id": { "view": true, "edit": false },
      "project_id": { "view": true, "edit": false },
      "title": { "view": true, "edit": true },
      "description": { "view": true, "edit": true },
      "status": { "view": true, "edit": true },
      "assignees": { "view": true, "edit": true },
      "created_by": { "view": true, "edit": false },
      "started_at": { "view": true, "edit": true },
      "completed_at": { "view": true, "edit": true }
    }
  },
  "users": {
    "view": true,
    "create": false,
    "edit": false,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": false },
      "email": { "view": true, "edit": false },
      "role": { "view": true, "edit": false }
    }
  }
}'),

('EDITOR', '{
  "projects": {
    "view": true,
    "create": false,
    "edit": true,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": true },
      "description": { "view": true, "edit": true },
      "created_by": { "view": true, "edit": false }
    }
  },
  "tasks": {
    "view": true,
    "create": true,
    "edit": true,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "project_id": { "view": true, "edit": false },
      "title": { "view": true, "edit": true },
      "description": { "view": true, "edit": true },
      "status": { "view": true, "edit": true },
      "assignees": { "view": true, "edit": true },
      "created_by": { "view": false, "edit": false },
      "started_at": { "view": true, "edit": true },
      "completed_at": { "view": true, "edit": true }
    }
  },
  "users": {
    "view": true,
    "create": false,
    "edit": false,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": false },
      "email": { "view": true, "edit": false },
      "role": { "view": true, "edit": false }
    }
  }
}'),

('VIEWER', '{
  "projects": {
    "view": true,
    "create": false,
    "edit": false,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": false },
      "description": { "view": true, "edit": false },
      "created_by": { "view": false, "edit": false }
    }
  },
  "tasks": {
    "view": true,
    "create": false,
    "edit": false,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "project_id": { "view": true, "edit": false },
      "title": { "view": true, "edit": false },
      "description": { "view": true, "edit": false },
      "status": { "view": true, "edit": false },
      "assignees": { "view": true, "edit": false },
      "created_by": { "view": false, "edit": false },
      "started_at": { "view": true, "edit": false },
      "completed_at": { "view": true, "edit": false }
    }
  },
  "users": {
    "view": true,
    "create": false,
    "edit": false,
    "delete": false,
    "fields": {
      "id": { "view": true, "edit": false },
      "name": { "view": true, "edit": false },
      "email": { "view": true, "edit": false },
      "role": { "view": true, "edit": false }
    }
  }
}');
