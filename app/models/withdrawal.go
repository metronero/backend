package models

import "time"

type Withdrawal struct {
	Id           string    `json:"withdrawal_id"`
	MerchantName string    `json:"merchant_name"`
	AccountId    string    `json:"merchant_id,omitempty"`
	Amount       uint64    `json:"amount"`
	AmountFloat  string    `json:"amount_float,omitempty"`
	Date         time.Time `json:"date"`
}

type WithdrawalRequest struct {
	Address string `json:"address"`
}
