package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/pkg/apierror"
)

// Returns a list of withdrawals.
func GetAdminWithdrawal(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context())
	if err != nil {
		writeError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
