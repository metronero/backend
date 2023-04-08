package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/app/controllers"
	"gitlab.com/metronero/backend/utils/config"
)

func registerRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewareServerHeader)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", controllers.Login)
	r.Post("/register", controllers.Register)
	// r.Get("/health", controllers.Health)

	r.Route("/merchant", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Get("/", controllers.MerchantInfo)
			r.Post("/", controllers.MerchantUpdate)
			r.Route("/withdrawal", func(r chi.Router) {
				r.Get("/", controllers.MerchantGetWithdrawals)
				r.Post("/", controllers.MerchantWithdrawFunds)
			})
			r.Route("/payment", func(r chi.Router) {
				r.Get("/", controllers.MerchantGetPayments)
				r.Post("/", controllers.PostMerchantPayment)
			})
			r.Route("/template", func(r chi.Router) {
				r.Get("/", controllers.MerchantGetTemplate)
				r.Post("/", controllers.MerchantPostTemplate)
			})
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.JwtSecret))
			r.Use(jwtauth.Authenticator)
			r.Use(middlewareAdminArea)
			r.Get("/", controllers.AdminInfo)
			r.Route("/instance", func(r chi.Router) {
				r.Get("/", controllers.GetInstance)
				r.Post("/", controllers.EditInstance)
			})
			r.Route("/withdrawal", func(r chi.Router) {
				r.Get("/", controllers.AdminGetWithdrawals)
				r.Get("/{merchant_id}", controllers.GetWithdrawalsByMerchant)
			})
			r.Route("/payment", func(r chi.Router) {
				r.Get("/", controllers.AdminGetPayments)
				r.Get("/{merchant_id}", controllers.GetPaymentsByMerchant)
			})
			r.Route("/merchant", func(r chi.Router) {
				r.Get("/{merchant_id}", controllers.AdminGetMerchant)
				r.Post("/{merchant_id}", controllers.AdminUpdateMerchant)
				r.Delete("/{merchant_id}", controllers.AdminDeleteMerchant)
				r.Get("/", controllers.AdminGetMerchantList)
			})
		})
	})

	// For handling MoneroPay callbacks
	r.Post("/callback/{uuid}", controllers.Callback)

	return r
}
