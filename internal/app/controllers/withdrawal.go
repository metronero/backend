package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/moneropay"
	"gitlab.com/metronero/metronero-go/api"
	"gitlab.com/metronero/metronero-go/models"
)

// Returns a list of withdrawals.
func GetAdminWithdrawal(w http.ResponseWriter, r *http.Request) {
	p, err := queries.GetWithdrawals(r.Context())
	if err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(p)
}
