package controllers

import (
	"net/http"
	"encoding/json"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

// Recaps relevant activity to be displayed on the admin dashboard.
func AdminInfo(w http.ResponseWriter, r *http.Request) {
	info, err := queries.GetAdminInfo()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
	json.NewEncoder(w).Encode(info)
}
