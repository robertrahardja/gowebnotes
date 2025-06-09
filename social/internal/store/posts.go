package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	// Version   int       `json:"version"`
	// Comments  []Comment `json:"comments"`
	// User      User      `json:"user"`
}

// this is for Postgres
// other db types should have its own implementation
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	// QueryRowContext captures these values
	// so you can populate the Post struct with the database-generated ID and timestamps.
	err := s.db.QueryRowContext(
		// ctx: The context for request cancellation and timeouts
		// query: Your SQL INSERT statement with placeholders ($1, $2, etc.)
		// args: The values to substitute for placeholders:
		ctx,
		query,
		//
		// post.Content → $1
		// post.Title → $2
		// post.UserID → $3
		// pq.Array(post.Tags) → $4 (PostgreSQL array)
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		// 		QueryRowContext returns a *sql.Row object,
		// 	which is immediately chain with .Scan()
		// 	to extract the returned values into your struct fields.
		//
		// Key Characteristics
		//
		// Single Row: Unlike Query, it's designed for queries that return exactly one row
		// Context-Aware: Respects context cancellation and timeouts
		// Immediate Execution: The query executes immediately, not when you call Scan()
		// Error Handling: Errors are deferred until you call Scan()
		//
		// Why It's Perfect Here
		// The INSERT statement uses RETURNING id, created_at, updated_at,
		// which returns exactly one row with the newly created record's generated values.
		//
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
