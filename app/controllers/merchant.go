package controllers

import (
	"net/http"
	"encoding/json"

	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

// Recaps relevant activity to be displayed on the merchant dashboard.
func MerchantInfo(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid token")
	}

	id := token["id"].(string)
	info, err := queries.GetMerchantInfo(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
	json.NewEncoder(w).Encode(info)
}
