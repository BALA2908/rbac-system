package models

type Project struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	CreatedBy         string   `json:"created_by"`
	AssignedEmployees []string `json:"assigned_employees,omitempty"`
}
