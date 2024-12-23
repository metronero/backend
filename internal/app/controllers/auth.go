package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/jackc/pgx/v5"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/auth"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, err)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		helpers.WriteError(w, apierror.ErrRequired, nil)
		return
	}

	account, err := queries.UserLogin(r.Context(), creds.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			helpers.WriteError(w, apierror.ErrUnauthorized, err)
		} else {
			helpers.WriteError(w, apierror.ErrDatabase, err)
		}
		return
	}
	if err := auth.CompareHashAndPassword(account.PasswordHash, creds.Password); err != nil {
		helpers.WriteError(w, apierror.ErrUnauthorized, err)
		return
	}
	sess := session.GetSession(r)
	if err := sess.Set("username", creds.Username); err != nil {
		helpers.WriteError(w, apierror.ErrSession, err)
		return
	}
	if err := sess.Set("account_id", account.AccountId); err != nil {
		helpers.WriteError(w, apierror.ErrSession, err)
		return
	}
	if err := sess.Set("role", account.Role); err != nil {
		helpers.WriteError(w, apierror.ErrSession, err)
		return
	}
	json.NewEncoder(w).Encode(models.LoginResponse{Role: account.Role, AccountId: account.AccountId})
}

// Only the instance admin can register new users
func PostRegister(w http.ResponseWriter, r *http.Request) {
	var creds models.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		helpers.WriteError(w, apierror.ErrBadRequest, err)
		return
	}

	if creds.Username == "" || creds.Password == "" || creds.Role == "" {
		helpers.WriteError(w, apierror.ErrRequired, nil)
		return
	}

	passwordHashBytes, err := auth.HashPassword(creds.Password)
	if err != nil {
		helpers.WriteError(w, apierror.ErrHash, err)
		return
	}

	accId, err := queries.UserRegister(r.Context(), creds.Username,
		string(passwordHashBytes), creds.Role)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
		return
	}
	json.NewEncoder(w).Encode(models.CreateAccountResponse{AccountId: accId})

}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	sess.Destroy(w, r)
	w.WriteHeader(http.StatusOK)
}
