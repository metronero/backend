package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
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

func PostAdminWithdraw(w http.ResponseWriter, r *http.Request) {
	var req models.PostWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}
	var resp models.PostWithdrawResponse
	// TODO: support sweep_all once its added to Mpay
	tx, err := moneropay.WithdrawFunds(req.Address, req.Amount)
	if err != nil {
		helpers.WriteError(w, apierror.ErrMoneropay, err)
		return
	}
	resp.TxId = tx.TxHashList[0]
	resp.Amount = tx.Amount
	json.NewEncoder(w).Encode(resp)
}
