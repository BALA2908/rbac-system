DELETE FROM role_permissions;

INSERT INTO role_permissions (role, permissions) VALUES

('ADMIN', '{
  "projects": { "view": true, "create": true, "edit": true },
  "users":    { "view": true, "create": true, "edit": true }
}'),

('MANAGER', '{
  "projects": { "view": true, "create": true, "edit": true }
}'),

('EDITOR', '{
  "projects": { "view": true, "create": false, "edit": true }
}'),

('VIEWER', '{
  "projects": { "view": true, "create": false, "edit": false }
}');
