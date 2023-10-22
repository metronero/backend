package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/auth"
	"gitlab.com/metronero/metronero-go/api"
	"gitlab.com/metronero/metronero-go/models"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: instead of form values use json data
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

// Only the instance admin can register new users
func PostRegister(w http.ResponseWriter, r *http.Request) {
	var creds models.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, api.ErrBadRequest, err)
		return
	}

	if creds.Username == "" || creds.Password == "" || creds.Role == "" {
		writeError(w, api.ErrRequired, nil)
		return
	}

	passwordHashBytes, err := auth.HashPassword(creds.Password)
	if err != nil {
		writeError(w, api.ErrHash, err)
		return
	}

	if err := queries.UserRegister(r.Context(), creds.Username,
		string(passwordHashBytes), creds.Role); err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}

}

// Invalidates bearer token of the user. Stores it in invalid_tokens table in
// database until expiry of the token.
func PostLogout(w http.ResponseWriter, r *http.Request) {
	token := jwtauth.TokenFromHeader(r)
	if err := queries.InvalidateToken(r.Context(), token); err != nil {
		writeError(w, api.ErrDatabase, err)
		return
	}
}
