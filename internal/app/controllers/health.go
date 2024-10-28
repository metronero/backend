package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	instance, err := queries.InstanceVersion(r.Context())
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(instance)
}
