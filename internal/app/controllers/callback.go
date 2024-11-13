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
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: is this accurate?
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}
	var data model.ReceiveGetResponse
	if err := json.Unmarshal(b, &data); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, nil)
		return
	}

	id := chi.URLParam(r, "invoice_id")
	completeOn, err := queries.FindInvoiceOwnerCompleteOn(context.Background(), id)

	var status string
	if completeOn == 10 {
		if data.Complete {
			status = "Completed"
		} else if data.Amount.Expected <= data.Amount.Covered.Total {
			status = "Confirming"
		} else if data.Amount.Expected > data.Amount.Covered.Total {
			status = "Partial"
		}
	} else {
		if completeOn == 1 {
			allHaveConfirmation := true
			atLeastOne := false
			for _, tx := range data.Transactions {
				if tx.Confirmations >= 1 {
					atLeastOne = true
				} else {
					allHaveConfirmation = false
				}
			}
			if allHaveConfirmation && data.Amount.Covered.Total >= data.Amount.Expected {
				status = "Completed"
			} else if !allHaveConfirmation && atLeastOne {
				if data.Amount.Covered.Total == data.Amount.Expected {
					status = "Confirming"
				} else if data.Amount.Covered.Total < data.Amount.Expected {
					status = "Partial"
				}
			}
		} else {
			// 0-conf
			if data.Amount.Covered.Total >= data.Amount.Expected {
				status = "Completed"
			} else if data.Amount.Covered.Total < data.Amount.Expected && data.Amount.Covered.Total > 0 {
				status = "Partial"
			} else {
				status = "Pending"
			}
		}
	}

	callbackData := string(b)
	// TODO: return the balance from this function to check if overpay occurred
	if err := queries.UpdatePayment(context.Background(), id, status, callbackData, data.Amount.Covered.Total); err != nil {
		log.Error().Err(err).Msg("Failed to update payment status")
	}
	if data.Complete {
		go queries.UpdateBalances(context.Background(), id, data.Amount.Expected)
		// TODO: go utils.SendCallback(paymentId, callbackUrl, callbackData)
	}
}
