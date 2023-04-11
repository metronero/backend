package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"gitlab.com/metronero/backend/app/queries"
	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/utils/moneropay"
)

func AdminGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func MerchantWithdrawFunds(w http.ResponseWriter, r *http.Request) {
	var wr models.WithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&wr); err != nil {
		writeError(w, ErrBadRequest, err)
		return
	}

	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}

	id := token["id"].(string)
	name := token["username"].(string)
	amount, err := queries.GetWithdrawableAmount(r.Context(), id)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}

	if amount == 0 {
		writeError(w, ErrNoFunds, err)
		return
	}

	withdrawal := &models.Withdrawal{
		Id: uuid.New().String(),
		MerchantName: name,
		Amount: amount,
		Date: time.Now(),
		AccountId: id,
	}

	if err := queries.SaveWithdrawal(r.Context(), withdrawal); err != nil {
		writeError(w, ErrDatabase, err)
		return
	}

	if err := moneropay.WithdrawFunds(wr.Address, amount); err != nil {
		writeError(w, ErrWithdraw, err)
		return
	}

	json.NewEncoder(w).Encode(withdrawal)
}

func MerchantGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}
	id := token["id"].(string)
	p, err := queries.GetWithdrawalsByAccount(r.Context(), id)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func GetWithdrawalsByMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	p, err := queries.GetWithdrawalsByAccount(r.Context(), id)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
