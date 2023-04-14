package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/backend/pkg/models"
	"gitlab.com/metronero/backend/pkg/api"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
)

// Return all payments submitted by all merchants.
func GetAdminPayment(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetAllPayments(r.Context())
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Return all payments associated with the merchant ID.
func GetAdminPaymentById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "merchant_id")
	p, err := queries.GetPaymentsByAccount(r.Context(), id)
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Get payments belonging to the logged in merchant.
func GetMerchantPayment(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
		return
	}
	id := token["id"].(string)
	p, err := queries.GetPaymentsByAccount(r.Context(), id)
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// Create a new payment request.
func PostMerchantPayment(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
		return
	}
	merchantId := token["id"].(string)
	name := token["username"].(string)
	var req models.PostPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, api.ErrBadRequest, nil)
		return
	}
	paymentId := uuid.New().String()
	subaddress, err := moneropay.CreatePayment(req.Amount, paymentId)
	if err != nil {
		writeError(w, api.ErrMoneropay, nil)
		return
	}
	if err := queries.CreatePaymentRequest(r.Context(), paymentId, merchantId, name,
		subaddress, &req); err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	res := &models.PostPaymentResponse{PaymentId: paymentId, Address: subaddress}
	json.NewEncoder(w).Encode(res)
}

// Load merchant template, payment details and serve payment page.
func PaymentPageHandler(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "payment_id")
	// Get payment details
	p, err := queries.GetPaymentPageInfo(r.Context(), paymentId)
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}

	// Load template
	var t *template.Template
	t, err = template.ParseFiles("./data/merchant_templates/" + p.TemplateId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t, err = template.ParseFiles("./data/merchant_templates/default")
			if err != nil {
				writeError(w, api.ErrTemplateLoad, err)
				return
			}
		} else {
			writeError(w, api.ErrTemplateLoad, err)
			return
		}
	}

	png, err := qrcode.Encode(p.Address, qrcode.Medium, 256)
	if err != nil {
		writeError(w, api.ErrTemplateLoad, err)
	}
	p.Qr = base64.StdEncoding.EncodeToString(png)
	p.AmountFloat = walletrpc.XMRToDecimal(p.Amount)
	t.Execute(w, p)
}
