package handlers

import (
	repositories "rbac-backend/internal/repository"
)

type AdminHandler struct {
	UserRepo *repositories.UserRepository
}

func NewAdminHandler(repo *repositories.UserRepository) *AdminHandler {
	return &AdminHandler{UserRepo: repo}
}
