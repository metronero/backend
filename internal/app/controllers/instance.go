package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetAdminInstance(w http.ResponseWriter, r *http.Request) {
	instanceVer, err := queries.InstanceVersion(r.Context())
	if err != nil {
		writeError(w, apierror.ErrDatabase, err)
		return
	}

	mpayHealth, mpayVer, err := moneropay.CheckHealth()
	if err != nil {
		writeError(w, apierror.ErrMoneropay, err)
		return
	}

	json.NewEncoder(w).Encode(models.GetInstanceInfoResponse{
		Version:          instanceVer,
		MoneroPayVersion: mpayVer,
		Healthy:          mpayHealth, // TODO: more checks to evaluate health
	})
}
