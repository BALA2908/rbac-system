package repositories

import (
	"database/sql"

	"rbac-backend/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, name, email, password_hash, role, is_active)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive,
	)
	return err
}

func (r *UserRepository) ListUsers() ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, email, role, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var isActive int
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &isActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		u.IsActive = isActive == 1
		users = append(users, u)
	}
	return users, rows.Err()
}
