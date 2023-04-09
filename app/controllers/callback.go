package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gitlab.com/moneropay/moneropay/v2/pkg/model"

	"gitlab.com/metronero/backend/app/queries"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: is this accurate?
		writeError(w, ErrBadRequest, nil)
		return
	}
	var data model.ReceiveGetResponse
	if err := json.Unmarshal(b, &data); err != nil {
		writeError(w, ErrBadRequest, nil)
		return
	}

	var status string
	if data.Complete {
		status = "Completed"
	} else if data.Amount.Expected == data.Amount.Covered.Total {
		status = "Confirming"
	} else if data.Amount.Expected > data.Amount.Covered.Total {
		status = "Partial"
	}

	id := chi.URLParam(r, "payment_id")
	if err := queries.UpdatePayment(r.Context(), id, status, string(b)); err != nil {
		log.Error().Err(err).Msg("Failed to update payment status")
	}
}
