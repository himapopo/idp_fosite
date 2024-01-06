package route

import "github.com/go-chi/chi"

func Routes(r chi.Router) {
	
	r.Route("/oauth2", func(r chi.Router) {
		r.Get("/auth", authenticate)
		r.Post("/token", token)
		r.Post("/introspect", introspect)
		r.Get("/userinfo", userInfo)
		r.Post("/logout", logout)
		r.Post("/revoke", revoke)
	})
}
