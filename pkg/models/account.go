package models

import "time"

type Account struct {
	AccountId    string
	Username     string
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
