package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

// Update merchant account settings as a merchant.
func PostMerchant(w http.ResponseWriter, r *http.Request) {
	var settings models.MerchantSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	if err := queries.ConfigureMerchant(r.Context(), id, &settings); err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
	}
}

func GetMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	settings, err := queries.GetMerchantSettings(ctx, id)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
	}
	json.NewEncoder(w).Encode(settings)
}

func GetMerchantStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	list, err := queries.GetMerchantStats(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// Returns a list of all merchants for the administrator.
func GetAdminMerchant(w http.ResponseWriter, r *http.Request) {
	list, err := queries.GetMerchantList(r.Context())
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// Return information about a merchant as an administrator.
func GetAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	m, err := queries.GetMerchantById(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(m)
}

// Update merchant's setting as an administrator.
func PostAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	var settings models.AdminMerchantSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}
	if err := queries.AdminEditMerchant(r.Context(), &settings, id); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
	}
}

// Delete merchant account.
func DeleteAdminMerchantById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	if err := queries.DeleteMerchantById(r.Context(), id); err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
	}
}

func GetMerchantCount(w http.ResponseWriter, r *http.Request) {
	total, active, err := queries.GetMerchantCount(r.Context())
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(models.GetMerchantCountResponse{Count: total, Active: active})
}
