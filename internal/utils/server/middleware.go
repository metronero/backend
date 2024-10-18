package server

import (
	"net/http"

	"gitea.com/go-chi/session"
)

const version = "0.0.0"

func middlewareServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Metronero/"+version)
		next.ServeHTTP(w, r)
	})
}

func middlewareAdminArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		role := sess.Get("role")
		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareMerchantArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		role := sess.Get("role")
		roleStr, ok := role.(string)
		if !ok || roleStr != "merchant" {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
