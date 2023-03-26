package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ApiError struct {
	Msg string `json:"message"`
	Code int   `json:"code"`
	status int
}

var (
	ErrInvalidToken = &ApiError{Code: 1, Msg: "Invalid token.", status: 403}
	ErrUnauthorized = &ApiError{Code: 2, Msg: "Unknown username or password.", status: 401}
	ErrRequired = &ApiError{Code: 3, Msg: "Required field(s) can't be empty.", status: 400}
	ErrHash = &ApiError{Code: 4, Msg: "Failed to hash password.", status: 500}
	ErrTokenIssue = &ApiError{Code: 5, Msg: "Failed to issue token.", status: 500}
	ErrUserExists = &ApiError{Code: 6, Msg: "User already exists.", status: 400}
	ErrNoId = &ApiError{Code: 7, Msg: "Unknown resource ID.", status: 400}
	ErrBadRequest = &ApiError{Code: 8, Msg: "Invalid request body.", status: 400}
	ErrDatabase = &ApiError{Code: 10, Msg: "Database error.", status: 500}
)

func writeError(w http.ResponseWriter, apiErr *ApiError, err error) {
	w.WriteHeader(apiErr.status)
	json.NewEncoder(w).Encode(apiErr)
	log.Error().Err(err).Int("code", apiErr.Code).Msg(apiErr.Msg)
}
