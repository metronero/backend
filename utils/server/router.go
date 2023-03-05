package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/app/controllers"
	"gitlab.com/moneropay/metronero/metronero-backend/utils/config"
)

func registerRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewareServerHeader)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//r.Get("/health", controllers.Health)
	r.Post("/login", controllers.Login)
	r.Post("/register", controllers.Register)

	r.Route("/merchant", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Get("/", controllers.MerchantInfo)
			/*r.Post("/template", controllers.Template)
			r.Post("/payment", controllers.GetPayment)
			r.Post("/create_payment", controllers.CreatePayment)
			r.Post("/withdraw", controllers.Withdraw)*/
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Use(middlewareAdminArea)
			r.Get("/", controllers.AdminInfo)
		})
	})

	/*
	r.Get("/p/{uuid}", controllers.Payment)
	r.Post("/callback/{uuid}", controllers.Callback)
	*/

	return r
}
