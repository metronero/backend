package models

import "time"

type ApiTokenInfo struct {
	AccountId  string    `json:"account_id,omitempty"`
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"valid_until"`
}

type CallbackSecret struct {
	AccountId  string    `json:"account_id,omitempty"`
	SecretKey  string    `json:"secret_key"`
	ValidUntil time.Time `json:"valid_until"`
}
