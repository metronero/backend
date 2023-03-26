package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
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
