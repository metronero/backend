package models

type Merchant struct {
	AccountId string
	CommissionRate uint64
	WalletAddress string
	TemplateId string
	APIKeyId string
}

type MerchantStats struct {
	MerchantId string
	Balance int64
	LastLogin int64
	TotalCommissionPaid uint64
	TotalSales uint64
}

type MerchantAPIKeys struct {
	Id string
	Key string
	ValidUntil int64
}
