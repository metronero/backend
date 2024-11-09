package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/moneropay/go-monero/walletrpc"
)

// Load merchant template, payment details and serve payment page.
func PaymentPageHandler(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "invoice_id")
	ctx := r.Context()
	p, err := queries.GetPaymentPageInfo(ctx, paymentId)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}

	// When expired, purge from MoneroPay to stop receiving callbacks
	if p.Status != "Completed" && p.Status != "Expired" && time.Now().After(p.Expires) {
		p.Status = "Expired"
		go func() {
			if err := queries.MarkInvoiceExpired(context.Background(), paymentId); err != nil {
				log.Err(err).Str("invoice_id", p.InvoiceId).Msg("Failed to mark invoice as expired")
				return
			}
		}()
	}

	// Load template
	var t *template.Template
	t, err = template.ParseFiles("./data/merchant_templates/" + p.TemplateId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t, err = template.ParseFiles("./data/merchant_templates/default")
			if err != nil {
				helpers.WriteError(w, apierror.ErrTemplateLoad, err)
				return
			}
		} else {
			helpers.WriteError(w, apierror.ErrTemplateLoad, err)
			return
		}
	}

	png, err := qrcode.Encode(p.Address, qrcode.Medium, 256)
	if err != nil {
		helpers.WriteError(w, apierror.ErrTemplateLoad, err)
	}
	p.Qr = base64.StdEncoding.EncodeToString(png)
	p.AmountFloat = walletrpc.XMRToDecimal(p.Amount)
	t.Execute(w, p)
}

func PaymentPageJsonHandler(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "invoice_id")
	ctx := r.Context()
	p, err := queries.GetPaymentPageInfo(ctx, paymentId)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
