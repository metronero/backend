package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/metronero/backend/pkg/models"
	"gitlab.com/metronero/backend/pkg/api"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/auth"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, api.ErrRequired, nil)
		return
	}

	account, err := queries.UserLogin(r.Context(), username)
	if err := auth.CompareHashAndPassword(account.PasswordHash, password); err != nil {
		writeError(w, api.ErrUnauthorized, err)
		return
	}

	token, expiry, err := auth.CreateUserToken(username, account.AccountId, 1*time.Hour)
	if err != nil {
		writeError(w, api.ErrTokenIssue, err)
		return
	}

	json.NewEncoder(w).Encode(&models.ApiTokenInfo{Token: token, ValidUntil: expiry})
}

func PostRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, api.ErrRequired, nil)
		return
	}

	passwordHashBytes, err := auth.HashPassword(password)
	if err != nil {
		writeError(w, api.ErrHash, err)
		return
	}

	if err := queries.UserRegister(r.Context(), username, string(passwordHashBytes)); err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
}