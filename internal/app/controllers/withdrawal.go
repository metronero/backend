package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/metronero-go/api"
)

// Returns a list of withdrawals.
func GetAdminWithdrawal(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context())
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
