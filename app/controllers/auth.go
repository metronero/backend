package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/metronero/backend/app/models"
	"gitlab.com/metronero/backend/app/queries"
	"gitlab.com/metronero/backend/utils/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, ErrRequired, nil)
		return
	}

	account, err := queries.UserLogin(r.Context(), username)
	if err := auth.CompareHashAndPassword(account.PasswordHash, password); err != nil {
		writeError(w, ErrUnauthorized, err)
	}

	token, expiry, err := auth.CreateUserToken(username, account.AccountId, 1*time.Hour)
	if err != nil {
		writeError(w, ErrTokenIssue, err)
	}

	// TODO: spawn goroutine to update last login timestamp here

	json.NewEncoder(w).Encode(&models.ApiTokenInfo{Token: token, ValidUntil: expiry})
}

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, ErrRequired, nil)
		return
	}

	passwordHashBytes, err := auth.HashPassword(password)
	if err != nil {
		writeError(w, ErrHash, err)
	}

	if err := queries.UserRegister(r.Context(), username, string(passwordHashBytes)); err != nil {
		writeError(w, ErrDatabase, err)
		return
	}
}
