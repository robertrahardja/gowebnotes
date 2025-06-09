package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Password  password `json:"-"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	// IsActive  bool     `json:"is_active"`
	// RoleID    int64    `json:"role_id"`
	// Role      Role     `json:"role"`
}

// this is for Postgres
// other db types should have its own implementation
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email, role_id) VALUES 
    ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
    RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
