package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/pkg/api"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetAdminInstance(w http.ResponseWriter, r *http.Request) {
	instance, err := queries.InstanceInfo(r.Context())
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(instance)
}

func PostAdminInstance(w http.ResponseWriter, r *http.Request) {
	var conf models.Instance
	if err := json.NewDecoder(r.Body).Decode(&conf); err != nil {
		writeError(w, api.ErrBadRequest, err)
		return
	}
	if err := queries.UpdateInstance(r.Context(), &conf); err != nil {
		writeError(w, api.ErrDatabase, err)
	}
}
