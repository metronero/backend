package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/pkg/api"
)

// Recaps relevant activity to be displayed on the admin dashboard.
func GetAdmin(w http.ResponseWriter, r *http.Request) {
	info, err := queries.GetAdminInfo(r.Context())
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(info)
}
