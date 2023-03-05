package models

import "time"

type Account struct {
	Id string
	Username string
	PasswordHash string
	CreationDate time.Time
}
