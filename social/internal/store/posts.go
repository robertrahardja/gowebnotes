package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
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

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at,  updated_at, tags 
		FROM posts
		WHERE id = $1
	`
	// ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	// defer cancel()

	var post Post

	// 	QueryRowContext() executes the SQL query and expects exactly one row back.
	// 	It immediately calls Scan() to copy the column values from that row into the provided variables (&post.ID, &post.UserID, etc.).
	// If the query finds no rows, Scan() returns sql.ErrNoRows. If it finds a row, Scan() populates your variables with the column data. Any database errors during execution get returned as err.
	// It's a one-shot operation: run query → get single row → fill variables → done.
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(post.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	// ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	// defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	// query := `
	// 	UPDATE posts
	// 	SET title = $1, content = $2, version = version + 1
	// 	WHERE id = $3 AND version = $4
	// 	RETURNING version
	// `
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3 

	`

	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	// ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	// defer cancel()
	//
	// err := s.db.QueryRowContext(
	// 	ctx,
	// 	query,
	// 	post.Title,
	// 	post.Content,
	// 	post.ID,
	// 	post.Version,
	// ).Scan(&post.Version)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return ErrNotFound
	// 	default:
	// 		return err
	// 	}
	// }

	return nil
}
