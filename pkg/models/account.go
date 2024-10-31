package models

import "time"

type Account struct {
	AccountId    string
	Username     string
	Role         string
	PasswordHash string
}

type AccountStats struct {
	CreationDate time.Time
	LastLogin    time.Time
}

type AccountChange struct {
	AccountId   string `json:"account_id"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type NewAccount struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	AccountId string `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Role      string `json:"role"`
	AccountId string `json:"account_id"`
}
