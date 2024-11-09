package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

// Return all payments submitted by all merchants.
func GetAdminPayment(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetAllPayments(r.Context())
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Return all payments associated with the merchant ID.
func GetAdminPaymentById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	p, err := queries.GetPaymentsByAccount(r.Context(), id, 0)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Get payments belonging to the logged in merchant.
func GetMerchantInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	p, err := queries.GetPaymentsByAccount(ctx, id, 0)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func GetMerchantInvoiceRecent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	p, err := queries.GetPaymentsByAccount(ctx, id, 10)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Create a new payment request.
func PostMerchantInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	merchantId := ctx.Value("account_id").(string)
	var req models.PostInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}
	paymentId := uuid.New().String()
	subaddress, err := moneropay.CreatePayment(req.Amount, paymentId)
	if err != nil {
		helpers.WriteError(w, apierror.ErrMoneropay, nil)
		return
	}
	if err := queries.CreatePaymentRequest(r.Context(), paymentId, merchantId,
		subaddress, &req); err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	res := &models.PostInvoiceResponse{InvoiceId: paymentId, Address: subaddress}
	json.NewEncoder(w).Encode(res)
}

func GetInvoiceCount(w http.ResponseWriter, r *http.Request) {
	total, pending, err := queries.GetInvoiceCount(r.Context())
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(&models.GetInvoiceCountResponse{Count: total, Pending: pending})
}

func GetAdminInvoiceRecent(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetPaymentsByAccount(r.Context(), "", 10)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
