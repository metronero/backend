package models

type AdminDashboardInfo struct {
	Stats InstanceStats `json:"instance_stats"`
	Recent []Withdrawal `json:"recent_withdrawals"`
}
