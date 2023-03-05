package server

import "net/http"

const version = "0.0.0"

func middlewareServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Metronero/" + version)
		next.ServeHTTP(w, r)
	})
}
