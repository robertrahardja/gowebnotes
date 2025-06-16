package main

import (
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// func (app *application) mount() *chi.Mux {
func (app *application) mount() http.Handler {
	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome there"))
	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})
		
		// TODO: Implement user routes when handlers are ready
		// r.Route("/users", func(r chi.Router) {
		// 	r.Put("/activate/{token}", app.activateUserHandler)
		//
		// 	r.Route("/{userID}", func(r chi.Router) {
		// 		r.Use(app.AuthTokenMiddleware)
		//
		// 		r.Get("/", app.getUserHandler)
		// 		r.Put("/follow", app.followUserHandler)
		// 		r.Put("/unfollow", app.unfollowUserHandler)
		// 	})
		//
		// 	r.Group(func(r chi.Router) {
		// 		r.Use(app.AuthTokenMiddleware)
		// 		r.Get("/feed", app.getUserFeedHandler)
		// 	})
		// })

		// TODO: Implement authentication routes
		// r.Route("/authentication", func(r chi.Router) {
		// 	r.Post("/user", app.registerUserHandler)
		// 	r.Post("/token", app.createTokenHandler)
		// })
	})

	// return mux
	return r
}

// func (app *application) run(mux *chi.Mux) error {
func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
