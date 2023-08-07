package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gitlab.com/moneropay/moneropay/v2/pkg/model"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/metronero-go/api"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: is this accurate?
		writeError(w, api.ErrBadRequest, nil)
		return
	}
	var data model.ReceiveGetResponse
	if err := json.Unmarshal(b, &data); err != nil {
		writeError(w, api.ErrBadRequest, nil)
		return
	}

	var status string
	if data.Complete {
		status = "Completed"
	} else if data.Amount.Expected <= data.Amount.Covered.Total {
		status = "Confirming"
	} else if data.Amount.Expected > data.Amount.Covered.Total {
		status = "Partial"
	}

	id := chi.URLParam(r, "payment_id")
	callbackData := string(b)
	// TODO: return the balance from this function to check if overpay occurred
	if err := queries.UpdatePayment(context.Background(), id, status, callbackData); err != nil {
		log.Error().Err(err).Msg("Failed to update payment status")
	}
	if data.Complete {
		go queries.UpdateBalances(context.Background(), id, data.Amount.Expected)
		// TODO: go utils.SendCallback(paymentId, callbackUrl, callbackData)
	}
}
