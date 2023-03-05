package server

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

const version = "0.0.0"

func middlewareServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Metronero/" + version)
		next.ServeHTTP(w, r)
	})
}

func middlewareAdminArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, c, err := jwtauth.FromContext(r.Context())
		if c["username"].(string) != "admin" || err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
			    http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
