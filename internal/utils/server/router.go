package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/internal/app/controllers"
	"gitlab.com/metronero/backend/internal/utils/config"
)

func registerRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewareServerHeader)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", controllers.PostLogin)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.JwtSecret))
		r.Use(jwtauth.Authenticator)
		r.Post("/logout", controllers.PostLogout)
	})

	// r.Get("/health", controllers.Health)

	r.Route("/merchant", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Get("/", controllers.GetMerchant)
			r.Post("/", controllers.PostMerchant)
			r.Route("/payment", func(r chi.Router) {
				r.Get("/", controllers.GetMerchantPayment)
				r.Post("/", controllers.PostMerchantPayment)
			})
			r.Route("/template", func(r chi.Router) {
				r.Get("/", controllers.GetMerchantTemplate)
				r.Post("/", controllers.PostMerchantTemplate)
				r.Delete("/", controllers.DeleteMerchantTemplate)
			})
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Use(middlewareAdminArea)
			r.Post("/register", controllers.PostRegister)
			r.Route("/instance", func(r chi.Router) {
				r.Get("/", controllers.GetAdminInstance)
				r.Post("/", controllers.PostAdminInstance)
			})
			r.Route("/withdrawal", func(r chi.Router) {
				r.Get("/", controllers.GetAdminWithdrawal)
			})
			r.Route("/payment", func(r chi.Router) {
				r.Get("/", controllers.GetAdminPayment)
				r.Get("/{merchant_id}", controllers.GetAdminPaymentById)
			})
			r.Route("/merchant", func(r chi.Router) {
				r.Get("/{merchant_id}", controllers.GetAdminMerchantById)
				r.Post("/{merchant_id}", controllers.PostAdminMerchantById)
				r.Delete("/{merchant_id}", controllers.DeleteAdminMerchantById)
				r.Get("/", controllers.GetAdminMerchant)
			})
		})
	})

	// For handling MoneroPay callbacks
	r.Post("/callback/{invoice_id}", controllers.CallbackHandler)

	// For handling payment page requests
	r.Get("/p/{invoice_id}", controllers.PaymentPageHandler)

	return r
}
