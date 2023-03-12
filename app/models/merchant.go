package models

import "time"

type Merchant struct {
	AccountId string `json:"id"`
	Name string `json:"name"`
	CommissionRate uint64 `json:"commission_rate"`
	WalletAddress string `json:"wallet_address"`
	TemplateId string `json:"template_id"`
}

type MerchantStats struct {
	AccountId string `json:"merchant_id"`
	Balance uint64 `json:"balance"`
	TotalSales uint64 `json:"total_sales"`
}

type MerchantAPIKeys struct {
	AccountId string `json:"merchant_id"`
	KeyId string `json:"key_id"`
	Secret string `json:"secret"`
	ValidUntil time.Time `json:"valid_until"`
}

type MerchantDashboardInfo struct {
	Stats MerchantStats `json:"merchant_stats"`
	Recent []Payment `json:"recent_payments"`
}
