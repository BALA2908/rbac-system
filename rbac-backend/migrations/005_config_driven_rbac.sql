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
  }
}');
