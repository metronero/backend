package models

type AdminDashboardInfo struct {
	Stats             InstanceStats `json:"instance_stats"`
	RecentWithdrawals []Withdrawal  `json:"recent_withdrawals"`
}
