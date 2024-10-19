package models

type GetBalanceResponse struct {
	Total    uint64 `json:"total"`
	Unlocked uint64 `json:"unlocked"`
}
