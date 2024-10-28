package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
)

func ListApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := ctx.Value("account_id").(string)
	keys, err := queries.GetAccountApiKeys(ctx, userId)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(keys)
}

func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := ctx.Value("account_id").(string)
	resp, appErr := queries.CreateApiKey(ctx, userId)
	if appErr != nil {
		helpers.WriteError(w, apierror.ErrDatabase, appErr)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func RevokeApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountId, _ := ctx.Value("account_id").(string)
	if appErr, err := queries.RevokeApiKey(ctx, chi.URLParam(r, "keyID"), accountId); appErr != nil {
		helpers.WriteError(w, appErr, err)
		return
	}
}
