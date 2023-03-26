package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

// Recaps relevant activity to be displayed on the admin dashboard.
func AdminInfo(w http.ResponseWriter, r *http.Request) {
	info, err := queries.GetAdminInfo()
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(info)
}
