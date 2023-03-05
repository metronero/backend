package models

type Merchant struct {
	Id string
	Username string
	PasswordHash string
	CommissionRate uint64
	WalletAddress string
	TemplateId string
	WithdrawalId uint64
	APIKeyId string
}

type MerchantStats struct {
	MerchantId string
	LastLogin int64
	AccountCreation uint64
	TotalCommissionPaid uint64
	TotalSales uint64
}

type MerchantAPIKeys struct {
	Id string
	Key string
	ValidUntil int64
}
