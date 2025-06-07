package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type api struct {
	addr string
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var users = []User{}

// This handles POST requests to create new users. It:
// Decodes JSON from the request body into a User struct
// Creates a new user with the provided first and last names
// Calls insertUser() to validate and add the user to the collection
// Returns HTTP 201 (Created) on success, or 400 (Bad Request) on error
func (a *api) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	var payload User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}

	// users = append(users, u)
	if err = insertUser(u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// This handles GET requests to retrieve all users. It:
//
// Sets HTTP 200 (OK) status
// Encodes the users slice as JSON and writes it to the response
// Returns HTTP 500 (Internal Server Error) if JSON encoding fails
func (a *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// This is a validation and insertion function that:
//
// Validates that both FirstName and LastName are provided (not empty)
// Checks for duplicate users by comparing both names
// Appends the user to the global users slice if validation passes
// Returns appropriate error messages for validation failures
func insertUser(u User) error {
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	for _, user := range users {
		if user.FirstName == u.FirstName && user.LastName == u.LastName {
			return errors.New("User already exists")
		}
	}

	users = append(users, u)
	return nil
}

// This sets up and starts the HTTP server:
//
// Creates an api instance configured to listen on port 8080
// Sets up HTTP routing using http.NewServeMux()
// Registers two routes: GET /users and POST /users
// Starts the server with ListenAndServe(), panicking if it fails to start
//
// The program creates a basic CRUD API where you can POST to create users and GET to retrieve all users, with in-memory storage using a global slice.
func main() {
	api := &api{addr: ":8080"}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUsersHandler)

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
