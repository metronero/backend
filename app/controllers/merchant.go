package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

// Recaps relevant activity to be displayed on the merchant dashboard.
func MerchantInfo(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}
	id := token["id"].(string)
	info, err := queries.GetMerchantInfo(r.Context(), id)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func MerchantUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	var settings models.MerchantSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeError(w, ErrBadRequest, nil)
		return
	}
	if err := queries.ConfigureMerchant(r.Context(), id, &settings); err != nil {
		writeError(w, ErrDatabase, nil)
	}
}

func AdminGetMerchantList(w http.ResponseWriter, r *http.Request) {
	list, err := queries.GetMerchantList(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func AdminGetMerchant(w http.ResponseWriter, r *http.Request) {
}

func AdminUpdateMerchant(w http.ResponseWriter, r *http.Request) {
}

func AdminDeleteMerchant(w http.ResponseWriter, r *http.Request) {
}
