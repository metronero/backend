package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	total, unlocked, err := moneropay.GetBalance()
	if err != nil {
		helpers.WriteError(w, apierror.ErrMoneropay, err)
	}
	json.NewEncoder(w).Encode(&models.GetBalanceResponse{Total: total, Unlocked: unlocked})

}
