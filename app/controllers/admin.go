package controllers

import "net/http"

func AdminInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
