package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// JSON Parsing: It reads JSON data from the request body into a CreatePostPayload struct using a helper function readJSON().
	// If parsing fails, it returns a 400 Bad Request error.
	// UserID hardcoded to 1 (with a TODO comment indicating this needs to change when authentication is implemented)
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	// Post Creation: It creates a new store.Post struct with:
	//
	// Title, Content, and Tags from the payload
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after auth
		UserID: 1,
	}

	ctx := r.Context()

	// It calls app.store.Posts.Create() to save the post to the database.
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return

	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// from /{postID}
	idParam := chi.URLParam(r, "postID")

	// Base 10: Parses the string as a decimal number (not binary/hex)
	// 64-bit: Uses int64 type, which can handle very large numbers (up to 9,223,372,036,854,775,807)
	// This is future-proof - even if your app starts small, you won't hit ID limits as it grows
	//
	// Type Safety
	// The database method app.store.Posts.GetByID(ctx, id) likely expects an int64 parameter, not a string.
	// Go is strongly typed, so you must convert between types explicitly.
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Gets the request context (useful for timeouts/cancellation)
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		//
		// If it's a "not found" error, returns HTTP 404
		// For any other error, returns HTTP 500
		//
		//
		// Uses helper functions to write JSON error responses
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
