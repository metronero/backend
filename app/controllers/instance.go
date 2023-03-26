package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

func GetInstance(w http.ResponseWriter, r *http.Request) {
	instance, err := queries.InstanceInfo(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, nil)
		return
	}
        json.NewEncoder(w).Encode(instance)
}

func EditInstance(w http.ResponseWriter, r *http.Request) {}
