package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/app/queries"
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
	// TODO: Initiate transfer with MoneroPay here
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
