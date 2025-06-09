package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	// The time.ParseDuration(maxIdleTime) call is necessary because maxIdleTime comes in as a string parameter,
	// but db.SetConnMaxIdleTime() requires a time.Duration type.
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//This line creates a context with a timeout for controlling how long an operation can run before being cancelled.
	//
	// context.Background():
	// Creates a basic, empty context that serves as the root.
	// It's never cancelled and has no deadline.
	//
	// context.WithTimeout(..., 5*time.Second):
	// Wraps the background context with a 5-second timeout.
	//
	// This returns two things:
	// ctx: A new context that will automatically cancel after 5 seconds
	// cancel: A function you can call to cancel the context early
	//
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Creates a 5-second timeout context and  calls PingContext() to
	// verify the database is actually reachable and the connection string is valid.
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
