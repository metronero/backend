package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	mpayHealth, _, err := moneropay.CheckHealth()
	if err != nil {
		helpers.WriteError(w, apierror.ErrMoneropay, err)
		return
	}
	json.NewEncoder(w).Encode(models.GetHealthResponse{Healthy: mpayHealth})
}
