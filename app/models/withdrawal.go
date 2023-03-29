package models

import "time"

type Withdrawal struct {
	Id           string    `json:"withdrawal_id"`
	MerchantName string    `json:"merchant_name"`
	Amount       uint64    `json:"amount"`
	Date         time.Time `json:"date"`
}
