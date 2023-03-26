package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func AdminGetPayments(w http.ResponseWriter, r *http.Request) {
        p, err := queries.GetPayments(r.Context())
        if err != nil {
                writeError(w, ErrDatabase, err)
		return
        }
        json.NewEncoder(w).Encode(p)
}

func MerchantGetPayments(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
        if err != nil {
                writeError(w, ErrInvalidToken, err)
		return
        }
        id := token["id"].(string)
        p, err := queries.GetPaymentsByAccount(r.Context(), id)
        if err != nil {
                writeError(w, ErrDatabase, err)
		return
        }
        json.NewEncoder(w).Encode(p)
}

func GetPaymentsByMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
        p, err := queries.GetPaymentsByAccount(r.Context(), id)
        if err != nil {
                writeError(w, ErrDatabase, err)
		return
        }
        json.NewEncoder(w).Encode(p)
}

func MerchantCreatePaymentReq(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
        if err != nil {
                writeError(w, ErrInvalidToken, err)
		return
        }
        id := token["id"].(string)
        name := token["username"].(string)
	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrBadRequest, nil)
		return
	}
	res, err := queries.CreatePaymentRequest(r.Context(), id, name, &req)
	if err != nil {
		writeError(w, ErrDatabase, nil)
	}
	json.NewEncoder(w).Encode(res)
}
