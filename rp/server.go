package main

import (
	"idp_fosite/route"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		w.Write([]byte(string(body)))
		return
	})

	route.Routes(r)

	http.ListenAndServe(":3846", r)
}
