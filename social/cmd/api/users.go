package main

import (
	"net/http"
	"social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// TODO: Implement user handlers when required interfaces are ready

// // GetUser godoc
// //
// //	@Summary		Fetches a user profile
// //	@Description	Fetches a user profile by ID
// //	@Tags			users
// //	@Accept			json
// //	@Produce		json
// //	@Param			id	path		int	true	"User ID"
// //	@Success		200	{object}	store.User
// //	@Failure		400	{object}	error
// //	@Failure		404	{object}	error
// //	@Failure		500	{object}	error
// //	@Security		ApiKeyAuth
// //	@Router			/users/{id} [get]
// func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
// 	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
//
// 	user, err := app.getUser(r.Context(), userID)
// 	if err != nil {
// 		switch err {
// 		case store.ErrNotFound:
// 			app.notFoundResponse(w, r, err)
// 			return
// 		default:
// 			app.internalServerError(w, r, err)
// 			return
// 		}
// 	}
//
// 	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
// 		app.internalServerError(w, r, err)
// 	}
// }

// TODO: Implement follow/unfollow functionality
// // FollowUser godoc
// // UnfollowUser godoc  
// // ActivateUser godoc

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
