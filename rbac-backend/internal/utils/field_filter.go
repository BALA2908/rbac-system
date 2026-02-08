// internal/utils/field_filter.go
package utils

import "rbac-backend/internal/models"

// FilterFields returns only fields the role is allowed to view.
// If fieldPerms is nil or empty, allows all (used for ADMIN full access).
func FilterFields(
	data map[string]interface{},
	fieldPerms map[string]models.FieldPermission,
) map[string]interface{} {
	if len(fieldPerms) == 0 {
		result := make(map[string]interface{}, len(data))
		for k, v := range data {
			result[k] = v
		}
		return result
	}

	result := make(map[string]interface{})
	for field, value := range data {
		perm, exists := fieldPerms[field]
		if !exists || !perm.View {
			continue
		}
		result[field] = value
	}
	return result
}

// FilterEditableFields returns only fields the role is allowed to create/edit.
// If fieldPerms is nil or empty, allows all (used for ADMIN full access).
func FilterEditableFields(
	data map[string]interface{},
	fieldPerms map[string]models.FieldPermission,
) map[string]interface{} {
	if len(fieldPerms) == 0 {
		result := make(map[string]interface{}, len(data))
		for k, v := range data {
			result[k] = v
		}
		return result
	}

	result := make(map[string]interface{})
	for field, value := range data {
		perm, exists := fieldPerms[field]
		if !exists || !perm.Edit {
			continue
		}
		result[field] = value
	}
	return result
}
