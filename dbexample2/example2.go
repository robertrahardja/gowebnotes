package main

import (
	"database/sql"
	"fmt"
)

type application struct {
	store Store
}

type UserRepository interface {
	GetByID(id int) (*User, error)
}

type Store interface {
	GetByID(id int) (*User, error)
}

type User struct {
	ID   string
	Name string
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) GetByID(id int) (*User, error) {
	row := r.db.QueryRow("SELECT id, name FROM users WHERE id = $1", id)

	var user User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type InMemRepository struct {
	users []User
}

func (r *InMemRepository) GetByID(id int) (*User, error) {
	return nil, nil
}

func main() {
	connStr := "user=username dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// clear separation of layers
	userRepository := NewPostgresUserRepository(db)
	userService := NewUserService(userRepository)

	app := &application{
		store: UserRepository,
	}

	user, err := userService.GetUserByID(1)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("User: %+v\n", user)
	}
}
