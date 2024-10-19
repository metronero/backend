package controllers

import (
	"encoding/json"
	"net/http"

	"gitea.com/go-chi/session"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/auth"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, apierror.ErrBadRequest, err)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		writeError(w, apierror.ErrRequired, nil)
		return
	}

	account, err := queries.UserLogin(r.Context(), creds.Username)
	if err != nil {
		writeError(w, apierror.ErrDatabase, err)
		return
	}
	if err := auth.CompareHashAndPassword(account.PasswordHash, creds.Password); err != nil {
		writeError(w, apierror.ErrUnauthorized, err)
		return
	}
	sess := session.GetSession(r)
	if err := sess.Set("username", account.Username); err != nil {
		writeError(w, apierror.ErrSession, err)
		return
	}
	if err := sess.Set("accountid", account.AccountId); err != nil {
		writeError(w, apierror.ErrSession, err)
		return
	}
	if err := sess.Set("role", account.Role); err != nil {
		writeError(w, apierror.ErrSession, err)
		return
	}
	json.NewEncoder(w).Encode(models.LoginResponse{Role: account.Role})
}

// Only the instance admin can register new users
func PostRegister(w http.ResponseWriter, r *http.Request) {
	var creds models.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, apierror.ErrBadRequest, err)
		return
	}

	if creds.Username == "" || creds.Password == "" || creds.Role == "" {
		writeError(w, apierror.ErrRequired, nil)
		return
	}

	passwordHashBytes, err := auth.HashPassword(creds.Password)
	if err != nil {
		writeError(w, apierror.ErrHash, err)
		return
	}

	accId, err := queries.UserRegister(r.Context(), creds.Username,
		string(passwordHashBytes), creds.Role)
	if err != nil {
		writeError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(models.CreateAccountResponse{AccountId: accId})

}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	sess.Destroy(w, r)
	w.WriteHeader(http.StatusOK)
}
