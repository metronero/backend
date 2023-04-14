package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/app/queries"
)

// Recaps relevant activity to be displayed on the merchant dashboard.
// Returns recent payments and merchant's withdrawable balance, total sales.
func GetMerchant(w http.ResponseWriter, r *http.Request) {
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

// Update merchant account settings as a merchant.
func PostMerchant(w http.ResponseWriter, r *http.Request) {
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

// Returns a list of all merchants for the administrator.
func GetAdminMerchant(w http.ResponseWriter, r *http.Request) {
	list, err := queries.GetMerchantList(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// Return information about a merchant as an administrator.
func GetAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	m, err := queries.GetMerchantById(r.Context(), id)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(m)
}

// Update merchant's setting as an administrator.
func PostAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	var settings models.MerchantSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeError(w, ErrBadRequest, nil)
		return
	}
	if err := queries.AdminEditMerchant(r.Context(), &settings); err != nil {
		writeError(w, ErrBadRequest, nil)
	}
}

// Delete merchant account.
func DeleteAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	if err := queries.DeleteMerchantById(r.Context(), id); err != nil {
		writeError(w, ErrDatabase, err)
	}
}
