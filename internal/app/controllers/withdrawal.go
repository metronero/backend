package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
)

// Returns a list of withdrawals.
func GetAdminWithdrawal(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context(), 0)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func GetAdminWithdrawalRecent(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context(), 10)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
