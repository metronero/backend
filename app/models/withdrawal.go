package models

import "time"

type Withdrawal struct {
	MerchantName string
	Amount uint64
	Date time.Time
}
