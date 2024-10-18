package models

import "time"

type Merchant struct {
	AccountId      string `json:"id"`
	Name           string `json:"name"`
	CommissionRate uint64 `json:"commission_rate"`
	WalletAddress  string `json:"wallet_address"`
	TemplateId     string `json:"template_id"`
	Disabled       bool   `json:"disabled"`
}

type MerchantStats struct {
	AccountId  string `json:"merchant_id,omitempty"`
	Balance    uint64 `json:"balance"`
	TotalSales uint64 `json:"total_sales"`
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

type MerchantDashboardInfo struct {
	Stats          MerchantStats `json:"merchant_stats"`
	RecentPayments []Payment     `json:"recent_payments"`
}
