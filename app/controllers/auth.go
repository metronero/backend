package controllers

import (
	"time"
	"encoding/json"
	"net/http"

	"gitlab.com/moneropay/metronero/metronero-backend/utils/auth"
	"gitlab.com/moneropay/metronero/metronero-backend/app/queries"
)

type LoginResponse struct {
	Token string `json:"token"`
	ValidUntil string `json:"valid_until"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, http.StatusBadRequest, "Required field(s) can't be empty")
		return
	}

	account, err := queries.UserLogin(r.Context(), username)
	if err := auth.CompareHashAndPassword(account.PasswordHash, password); err != nil {
		writeError(w, http.StatusUnauthorized, "Unknown username or password")
	}

	token, expiry, err := auth.CreateUserToken(username, account.Id, 1 * time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to issue token")
	}

	// TODO: spawn goroutine to update last login timestamp here

	json.NewEncoder(w).Encode(&LoginResponse{Token: token, ValidUntil: expiry})
}

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeError(w, http.StatusBadRequest, "Required field(s) can't be empty")
		return
	}

	passwordHashBytes, err := auth.HashPassword(password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to hash password: " + err.Error())
	}

	if err := queries.UserRegister(r.Context(), username, string(passwordHashBytes)); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to register user: " + err.Error())
		return
	}
}
