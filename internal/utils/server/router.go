package server

import (
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.CorsOrigin},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowCredentials: true,
	}))

	r.Group(func(r chi.Router) {
		r.Use(session.Sessioner(session.Options{
			CookieName: "MNSession",
		}))

		r.Post("/login", controllers.PostLogin)

		r.Group(func(r chi.Router) {
			r.Use(middlewareAuthArea)
			r.Post("/logout", controllers.PostLogout)
			r.Get("/health", controllers.GetHealth)

			r.Route("/merchant", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(middlewareMerchantArea)
					r.Post("/", controllers.PostMerchant)
					r.Get("/stats", controllers.GetMerchantStats)
					r.Get("/fiatrate/{fiat}", controllers.GetFiatRate)
					r.Route("/invoice", func(r chi.Router) {
						r.Get("/recent", controllers.GetMerchantInvoiceRecent)
						r.Get("/", controllers.GetMerchantInvoice)
						r.Post("/", controllers.PostMerchantInvoice)
					})
					r.Route("/template", func(r chi.Router) {
						r.Get("/", controllers.GetMerchantTemplate)
						r.Get("/default", controllers.GetDefaultTemplatePreview)
						r.Post("/", controllers.PostMerchantTemplate)
						r.Delete("/", controllers.DeleteMerchantTemplate)
					})
					r.Route("/apikey", func(r chi.Router) {
						r.Get("/", controllers.ListApiKey)
						r.Post("/", controllers.CreateApiKey)
						r.Delete("/{keyID}", controllers.RevokeApiKey)
					})
				})
			})

			r.Route("/admin", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(middlewareAdminArea)
					r.Post("/register", controllers.PostRegister)
					r.Get("/instance", controllers.GetAdminInstance)
					r.Get("/balance", controllers.GetBalance)
					r.Route("/withdrawal", func(r chi.Router) {
						r.Get("/", controllers.GetAdminWithdrawal)
						r.Get("/recent", controllers.GetAdminWithdrawalRecent)
						r.Post("/", controllers.PostAdminWithdraw)
					})
					r.Route("/invoice", func(r chi.Router) {
						r.Get("/", controllers.GetAdminPayment)
						r.Get("/{merchant_id}", controllers.GetAdminPaymentById)
						r.Get("/count", controllers.GetInvoiceCount)
						r.Get("/recent", controllers.GetAdminInvoiceRecent)
					})
					r.Route("/merchant", func(r chi.Router) {
						r.Get("/{merchant_id}", controllers.GetAdminMerchantById)
						r.Post("/{merchant_id}", controllers.PostAdminMerchantById)
						r.Delete("/{merchant_id}", controllers.DeleteAdminMerchantById)
						r.Get("/", controllers.GetAdminMerchant)
						r.Get("/count", controllers.GetMerchantCount)
					})
				})
			})
		})
	})

	// For handling MoneroPay callbacks
	r.Post("/callback/{invoice_id}", controllers.CallbackHandler)

	assetServer := http.FileServer(http.Dir(config.TemplateDir))
	r.Handle("/merchantdata/*", http.StripPrefix("/merchantdata/", disableDirListing(assetServer)))

	// For handling payment page requests
	r.Get("/p/page/{invoice_id}", controllers.PaymentPageHandler)
	r.Get("/p/json/{invoice_id}", controllers.PaymentPageJsonHandler)

	r.Route("/merchant_api", func(r chi.Router) {
		r.Use(middlewareAdminArea)
		r.Post("/", controllers.PostMerchantInvoice)
	})

	return r
}
