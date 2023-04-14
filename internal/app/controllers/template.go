package controllers

import (
	"encoding/base64"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	qrcode "github.com/skip2/go-qrcode"

	"gitlab.com/metronero/backend/pkg/models"
	"gitlab.com/metronero/backend/pkg/api"
)

// Get a preview of the payment page template
func GetMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
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
				writeError(w, api.ErrTemplateLoad, err)
				return
			}
		} else {
			writeError(w, api.ErrTemplateLoad, err)
			return
		}
	}

	address := "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw"
	png, err := qrcode.Encode(address, qrcode.Medium, 256)
	if err != nil {
		writeError(w, api.ErrTemplateLoad, err)
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

// Upload new template.
func PostMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	// 20 MB max upload size
	// TODO: make this configurable
	r.ParseMultipartForm(20 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, api.ErrBadRequest, err)
		return
	}
	defer file.Close()

	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
		return
	}
	id := token["id"].(string)
	f, err := os.OpenFile("./data/merchant_templates/"+id, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		writeError(w, api.ErrTemplateSave, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

// Reset template back to default. Works by deleting merchant's template file.
func DeleteMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
		return
	}
	accountId := token["id"].(string)
	if err := os.Remove("./data/merchant_templates/" + accountId); err != nil && err != os.ErrNotExist {
		writeError(w, api.ErrTemplateDelete, err)
		return
	}
}
