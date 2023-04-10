package models

type Instance struct {
	Version           string `json:"version,omitempty"`
	DefaultCommission uint64 `json:"default_commission"`
	DefaultCommissionFloat string

	// Allow multiple merchants to use this instance and store their funds
	// in the instance's wallet.
	CustodialMode bool `json:"custodial_mode"`

	// If false, the instance is invite-only. If custodial mode is off
	// registrations can't be allowed.
	RegistrationsAllowed bool `json:"registrations_allowed"`

	// Decide when auto withdrawals should happen if custodial mode is on.
	WithdrawalTimes string `json:"withdrawal_times"`
}

type InstanceBootstrap struct {
	// Flag to indicate bootstrapping is needed:
	// Creation of admin account, default theme etc.
	FirstRun bool `json:"first_run"`
}

type InstanceStats struct {
	// Total fees paid by merchants
	WalletBalance      uint64 `json:"wallet_balance"`
	WalletBalanceFloat string
	TotalProfits       uint64 `json:"total_profits"`
	TotalProfitsFloat  string
	TotalMerchants     uint64 `json:"total_merchants"`
}

type InstanceInfo struct {
	Details Instance      `json:"details"`
	Stats   InstanceStats `json:"stats"`
}
