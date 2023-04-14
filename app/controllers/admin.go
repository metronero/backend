package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/app/queries"
)

// Recaps relevant activity to be displayed on the admin dashboard.
func GetAdmin(w http.ResponseWriter, r *http.Request) {
	info, err := queries.GetAdminInfo(r.Context())
	if err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(info)
}
