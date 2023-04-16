package controllers

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/internal/utils/config"
	"gitlab.com/metronero/backend/pkg/api"
	"gitlab.com/metronero/backend/pkg/models"
)

// Get a preview of the payment page template
func GetMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	_, token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		writeError(w, api.ErrInvalidToken, err)
		return
	}

	accountId := token["id"].(string)
	var t *template.Template
	p, err := url.JoinPath(config.TemplateDir, accountId)
	if err != nil {
		writeError(w, api.ErrTemplateLoad, err)
		return
	}
	t, err = template.ParseFiles(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t = config.DefaultPaymentTemplate
		} else {
			writeError(w, api.ErrTemplateLoad, err)
			return
		}
	}

	// Execute with dummy data
	t.Execute(w, &models.PaymentPageInfo{
		InvoiceId:    "0b1b6a94-9ec2-4bdf-8251-6a46aca6a332",
		MerchantName: token["username"].(string),
		Amount:       120000000000,
		AmountFloat:  "0.12",
		OrderId:      "AI6X21",
		Status:       "Pending",
		LastUpdate:   time.Now(),
		Address:      "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw",
		ExtraData:    "Sneakers x 2, Jacket x 1",
		Qr:           "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACuElEQVR42uyZwY30KhCEy+LQR0IgE5PYaOyRE8OZEAJHDhb11M386519ARhLw2nl/Q4Dbrqq2viu77rhEpLcy8QXmR8AaqwzJn1abgRkQBI6sGSEqshUADcSQErC7F9cgaC0NNnLwm0wYGey010DMwmMCZSJDY5rOIAoe51uB+iumHzDIxwBQZJomX9W1NVAv3qzf5VH2PIRJNVZ9v/dzUsBWxP3fvWoR00y/Wl1FwOSEYUllqdfsVBr+t3CzndxPUBJNRarUh569VKdKnzzvBGgS2vDvxTJLteonbkAy0CACsTsd+sP5Ls/6D/cSIAkbcRlMuVV8SDr7Nu5iwEAkJh9wlIe0KtXI/c6sfntRoDQCiay+QPO3gVmScBQQFZXo3rhVZqPgMgm6XwXQwBkkgaYyXmAWY+aDRgKyNp7y2wFvgb1tFBPqyc/DtB3MXH1zC67TFYIy3LqxQ0ABM0Lfi9Pv3HFEVRA/rTiywGh2ZykGhaYD6jnYVNjhnGAjKg/2rM4C2JMNVZtF2cQuxzQJ2J/e3LLzEJNi+oXbgUAMGn2zND4ADbE8mk4LwYkv6UZ5ic1UFZIw/PsciMA+sTvdrLvVqwRuPkVAwG15xwWx43OANm5/lbe6wH2Kn3C5X7UwjrxVdydAOtxGs/Kkh1dRjQftv5SvesB7XIaadQyaDS3LqeRZx0HQBcIdWJHOIJjjX1Q86O8gwD/7GJe6HIfeZ31cAtAiChN/QOWbnohTauo/Gj39UCfgRSL5hnQXeirwdN/DGquBd7TJIUQNv5MF59nCBoAsGGyZ3niAUcVD2sYSzmd2ACATThtF3DZsYI7+Ttf3AVQLSaPsMLexVTxZ5tDABZx+zSJ1HSJqXx8XLga6F9AtPfaNEkoezWNGwj49y1JNWzhEUwvtMwfuA/wXd811PovAAD//y9+6+A69LGyAAAAAElFTkSuQmCC",
	})
}

// Upload new template.
func PostMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << config.TemplateMaxSize)
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
	p, err := url.JoinPath(config.TemplateDir, id)
	if err != nil {
		writeError(w, api.ErrTemplateSave, err)
		return
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0644)
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
	p, err := url.JoinPath(config.TemplateDir, accountId)
	if err != nil {
		writeError(w, api.ErrTemplateDelete, err)
		return
	}
	if err := os.Remove(p); err != nil && err != os.ErrNotExist {
		writeError(w, api.ErrTemplateDelete, err)
		return
	}
}
