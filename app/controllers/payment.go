package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/app/queries"
	"gitlab.com/metronero/backend/utils/moneropay"
)

func AdminGetPayments(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetAllPayments(r.Context())
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

// Create a new payment request
func PostMerchantPayment(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}
	merchantId := token["id"].(string)
	name := token["username"].(string)
	var req models.PostPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrBadRequest, nil)
		return
	}
	paymentId := uuid.New().String()
	subaddress, err := moneropay.CreatePayment(req.Amount, paymentId)
	if err != nil {
		writeError(w, ErrMoneropay, nil)
		return
	}
	if err := queries.CreatePaymentRequest(r.Context(), paymentId, merchantId, name,
	    subaddress, &req); err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	res := &models.PostPaymentResponse{PaymentId: paymentId, Address: subaddress}
	json.NewEncoder(w).Encode(res)
}
