package models

import "time"

type Withdrawal struct {
	Id      string    `json:"withdrawal_id"`
	Amount  uint64    `json:"amount"`
	Address string    `json:"address"`
	Date    time.Time `json:"date"`
}

type PostWithdrawRequest struct {
	Address  string `json:"address"`
	Amount   uint64 `json:"amount,omitempty"`
	SweepAll bool   `json:"sweep_all,omitempty"`
}

type PostWithdrawResponse struct {
	Id     string `json:"withdrawal_id"`
	TxId   string `json:"txid"`
	Amount uint64 `json:"amount"`
}
