package controllers

import "net/http"

func MerchantInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
