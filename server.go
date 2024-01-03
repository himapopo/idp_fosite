package main

import (
	"net/http"

	"idp_fosite/route"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("welcome"))
		return
	})

	route.Routes(r)

	http.ListenAndServe(":3000", r)
}
