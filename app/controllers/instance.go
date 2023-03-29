package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/app/queries"
)

func GetInstance(w http.ResponseWriter, r *http.Request) {
	instance, err := queries.InstanceInfo(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(instance)
}

func EditInstance(w http.ResponseWriter, r *http.Request) {
	var conf models.Instance
	if err := json.NewDecoder(r.Body).Decode(&conf); err != nil {
		writeError(w, ErrBadRequest, err)
		return
	}
	if err := queries.UpdateInstance(r.Context(), &conf); err != nil {
		writeError(w, ErrDatabase, err)
	}
}
