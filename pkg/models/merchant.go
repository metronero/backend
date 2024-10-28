package models

import "time"

type Merchant struct {
	AccountId string `json:"id"`
	Name      string `json:"name"`
	Disabled  bool   `json:"disabled"`
}

type MerchantStats struct {
	TotalInvoices string `json:"total_invoices"`
	Pending       uint64 `json:"pending"`
	TotalSales    uint64 `json:"total_sales"`
}

type MerchantAPIKeys struct {
	AccountId  string    `json:"merchant_id"`
	KeyId      string    `json:"key_id"`
	Secret     string    `json:"secret"`
	ValidUntil time.Time `json:"valid_until"`
}

type MerchantSettings struct {
	AccountId      string  `json:"account_id,omitempty"`
	CommissionRate *uint64 `json:"commission_rate,omitempty"`
	Disabled       *bool   `json:"disabled,omitempty"`
}

type GetMerchantCountResponse struct {
	Count  uint64 `json:"count"`
	Active uint64 `json:"active"`
}
