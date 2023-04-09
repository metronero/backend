package controllers

import (
	"errors"
	"encoding/base64"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	qrcode "github.com/skip2/go-qrcode"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/app/queries"
)

// For previewing payment page template
func MerchantGetTemplate(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}

	accountId := token["id"].(string)
	name := token["username"].(string)
	var t *template.Template
	t, err = template.ParseFiles("./data/merchant_templates/" + accountId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t, err = template.ParseFiles("./data/merchant_templates/default")
			if err != nil {
				writeError(w, ErrTemplateLoad, err)
				return
			}
		} else {
			writeError(w, ErrTemplateLoad, err)
			return
		}
	}

	address := "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw"
	png, err := qrcode.Encode(address, qrcode.Medium, 256)
	if err != nil {
		writeError(w, ErrTemplateLoad, err)
	}

	// Execute with dummy data
	t.Execute(w, &models.PaymentPageInfo{
		InvoiceId:    "0b1b6a94-9ec2-4bdf-8251-6a46aca6a332",
		MerchantName: name,
		Amount:       120000000000,
		AmountFloat:  "0.12",
		OrderId:      "AI6X21",
		Status:       "Pending",
		LastUpdate:   time.Now(),
		Address:      address,
		ExtraData:    "Sneakers x 2, Jacket x 1",
		Qr:           base64.StdEncoding.EncodeToString(png),
	})
}

func MerchantPostTemplate(w http.ResponseWriter, r *http.Request) {
	// 20 MB max upload size
	// TODO: make this configurable
	r.ParseMultipartForm(20 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, ErrBadRequest, err)
		return
	}
	defer file.Close()

	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}
	id := token["id"].(string)
	f, err := os.OpenFile("./data/merchant_templates/"+id, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		writeError(w, ErrTemplateSave, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func GetPaymentPage(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, ErrInvalidToken, err)
		return
	}

	paymentId := chi.URLParam(r, "payment_id")
	// Get payment details
	p, err := queries.GetPaymentPageInfo(r.Context(), paymentId)
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}

	// Load template
	accountId := token["id"].(string)
	var t *template.Template
	t, err = template.ParseFiles("./data/merchant_templates/" + accountId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t, err = template.ParseFiles("./data/merchant_templates/default")
			if err != nil {
				writeError(w, ErrTemplateLoad, err)
				return
			}
		} else {
			writeError(w, ErrTemplateLoad, err)
			return
		}
	}

        png, err := qrcode.Encode(p.Address, qrcode.Medium, 256)
        if err != nil {
                writeError(w, ErrTemplateLoad, err)
        }
	p.Qr = base64.StdEncoding.EncodeToString(png)
	p.AmountFloat = walletrpc.XMRToDecimal(p.Amount)
	t.Execute(w, p)
}
