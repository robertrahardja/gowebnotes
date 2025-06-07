package main

import "net/http"

type api struct {
	addr string
}

func (s *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {

		case "/":
			w.Write([]byte("GET method"))

		case "/index":
			w.Write([]byte("GET index"))

		}
	case http.MethodPost:
		w.Write([]byte("POST method"))
	}
}

func main() {
	api := &api{addr: ":8080"}

	srv := &http.Server{
		Addr:    api.addr,
		Handler: api,
	}
	srv.ListenAndServe()
}
