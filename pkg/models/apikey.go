package models

import "time"

type ApiKey struct {
	KeyId  string    `json:"key_id"`
	Expiry time.Time `json:"expiry"`
}

type CreateApiKeyResponse struct {
	KeyId  string    `json:"key_id"`
	Key    string    `json:"key"`
	Expiry time.Time `json:"expiry"`
}
